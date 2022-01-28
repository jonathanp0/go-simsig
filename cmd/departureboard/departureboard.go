package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/go-stomp/stomp/v3"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/jonathanp0/go-simsig/gateway"
	"github.com/jonathanp0/go-simsig/wttxml"
)

// type definitions
type Location struct {
	Name string
	Code string
}

type LocationStopListMap map[string]*LocationStopList

// constants
var movementQueueName = "/topic/TRAIN_MVT_ALL_TOC"
var simsigQueueName = "/topic/SimSig"

//global variables
var locations []Location
var stopsAtLocations map[string]*LocationStopList
var stop = make(chan bool)
var currentClock gateway.ClockMsg

func main() {

	runtime.LockOSThread()
	runWindowsUI()

}

//Gateway Message Processing

//Process train_location gateway message
func processLocationMessage(m *gateway.TrainLocation, locations LocationStopListMap) {
	stops, ok := locations[m.Location]

	if !ok || m.Action == "pass" {
		//location with no scheduled calls, ignore
		return
	}

	/*if *verbose {
		println(prettyPrint(*m))
	}*/

	needsUpdate := false
	var stopTime uint

	//Search the location for this headcode and update
	for i, _ := range stops.Stops {
		stop := &stops.Stops[i]
		if stop.Headcode == m.Headcode {
			if m.Action == "arrive" {
				stop.Arrived = true
			} else if m.Action == "depart" {
				stop.Departed = true
			}
			stop.ActualPlatform = m.Platform
			if !stop.Updated {
				needsUpdate = true
				stopTime = stop.Time()
			}
		}
	}

	// Update all future stops to show this train is in the sim
	if needsUpdate {
		for _, location := range locations {
			for i, _ := range location.Stops {
				stop := &location.Stops[i]
				if (stop.Headcode == m.Headcode) && (stop.Time() > stopTime) {
					stop.Updated = true
				}
			}
		}
	}

}

//Process train_delay gateway message
func processDelayMessage(m *gateway.TrainDelay, locations LocationStopListMap) {

	/*if *verbose {
		println(prettyPrint(*m))
	}*/

	for _, location := range locations {
		for i, _ := range location.Stops {
			stop := &location.Stops[i]
			if stop.Headcode == m.Headcode {
				stop.Updated = true
				stop.DelaySeconds = m.Delay
			}
		}
	}
}

func gatewayConnection(user string, password string, address string) {
	//Iniate STOMP Connection
	subscribed := make(chan bool)

	go recvMessages(&currentClock, stopsAtLocations, subscribed, user, password, address)

	// wait until we know the receiver has subscribed
	<-subscribed

	go webInterface(locations, stopsAtLocations, &currentClock)
	webInterfaceReady()

	//run indefinitely
	<-stop
}

//Main communication thread for Interface Gateway
func recvMessages(clock *gateway.ClockMsg, locations LocationStopListMap, subscribed chan bool, user string, pass string, serverAddr string) {
	defer func() {
		stop <- true
	}()

	//login credentials
	var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{}

	if user != "" {
		options = append(options, stomp.ConnOpt.Login(user, pass))
	}

	conn, err := stomp.Dial("tcp", serverAddr, options...)

	if err != nil {
		updateStatus("cannot connect to server: " + err.Error())
		return
	}

	subMvt, err := conn.Subscribe(movementQueueName, stomp.AckAuto)
	if err != nil {
		updateStatus("cannot subscribe to " + movementQueueName + ": " + err.Error())
		return
	}
	subSimsig, err := conn.Subscribe(simsigQueueName, stomp.AckAuto)
	if err != nil {
		updateStatus("cannot subscribe to " + simsigQueueName + ": " + err.Error())
		return
	}
	conn.Send("/topic/SimSig", "text/plain", []byte("{\"idrequest\":{}}"))
	close(subscribed)

	//Wait for a message from either subscription
	for {
		select {
		case msg := <-subMvt.C:
			var decodedMsg gateway.TrainMovementMessage
			err := json.Unmarshal(msg.Body, &decodedMsg)
			if err != nil {
				showError("Error parsing Train Movement message: " + err.Error())
				continue
			}
			if decodedMsg.TrainLocation != nil {
				processLocationMessage(decodedMsg.TrainLocation, locations)
			} else if decodedMsg.TrainDelay != nil {
				processDelayMessage(decodedMsg.TrainDelay, locations)
			}
		case msg := <-subSimsig.C:
			var decodedMsg gateway.SimSigMessage
			err := json.Unmarshal(msg.Body, &decodedMsg)
			if err != nil {
				showError("Error parsing SimSig message: " + err.Error())
				continue
			}
			if decodedMsg.Clock != nil {
				*clock = *decodedMsg.Clock
			}
		}

	}
}

// Web Interface
func formatOptionalTime(seconds *uint) string {
	if seconds == nil {
		return "__:__"
	} else {
		return formatTime(*seconds)
	}
}

func formatTime(seconds uint) string {
	return fmt.Sprintf("%02d:%02d", seconds/3600, seconds%3600/60)
}

func localTemplatePath(tmpl string) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	return filepath.Join(exPath, tmpl)
}

func serveDepartureBoard(currentClock *gateway.ClockMsg, location string, stopList *LocationStopList, w http.ResponseWriter, req *http.Request) {
	var data struct {
		Clock    string
		ClockRaw uint
		Area     string
		Location string
		Stops    []LocationStop
	}
	data.Clock = formatTime(currentClock.Clock)
	data.ClockRaw = currentClock.Clock
	data.Area = currentClock.AreaID
	data.Location = location
	data.Stops = stopList.Stops

	tmpl, err := template.ParseFiles(localTemplatePath("tmpl/board.tmpl"))
	if err != nil {
		showError("board.tmpl error: " + err.Error())
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		showError("board.tmpl error: " + err.Error())
	}
}

//generate a handler function for a location
func serveDepartureBoardForLocation(currentClock *gateway.ClockMsg, location string, stopList *LocationStopList) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		serveDepartureBoard(currentClock, location, stopList, w, req)
	}
}

//web interface main loop
func webInterface(locations []Location, locationStops map[string]*LocationStopList, currentClock *gateway.ClockMsg) {

	for _, loc := range locations {
		http.HandleFunc("/board/"+loc.Code, serveDepartureBoardForLocation(currentClock, loc.Name, locationStops[loc.Code]))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Clock     string
			Area      string
			Locations []Location
		}
		data.Clock = formatTime(currentClock.Clock)
		data.Area = currentClock.AreaID
		data.Locations = locations

		tmpl, err := template.ParseFiles(localTemplatePath("tmpl/index.tmpl"))
		if err != nil {
			showError("index.tmpl error: " + err.Error())
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			showError("index.tmpl error: " + err.Error())
		}
	})

	http.ListenAndServe(":8090", nil)
}

// Timetable Parsing
func loadTimetable(filename string) string {

	//Read WTT
	data, err := wttxml.ReadSavedTimetable(filename)
	if err != nil {
		return ("Error reading WTT: " + err.Error())
	}

	//Build stop list from WTT
	var wtt wttxml.SimSigTimetable
	err = xml.Unmarshal(data, &wtt)
	if err != nil {
		return ("WTT Parsing Error: " + err.Error())
	}

	locationCodes := buildSortedLocationList(wtt.Timetables.Timetable)

	stopsAtLocations = buildLocationStopList(locationCodes, wtt.Timetables.Timetable)

	//Sort list by time and prune location list
	n := 0
	for name, locStops := range stopsAtLocations {
		sort.Sort(locStops)
		if locStops.Len() > 0 {
			locationCodes[n] = name
			n++
		}
	}

	locationCodes = locationCodes[:n]
	sort.Strings(locationCodes)

	locations = make([]Location, 0, len(locationCodes))

	//Attempt TIPLOC Lookups
	tiplocs := readTIPLOCs()

	for i, _ := range locationCodes {
		if desc, ok := tiplocs[locationCodes[i]]; ok {
			locations = append(locations, Location{strings.Title(strings.ToLower(desc)), locationCodes[i]})
		} else {
			locations = append(locations, Location{locationCodes[i], locationCodes[i]})
		}
	}

	sort.Slice(locations, func(i, j int) bool { return locations[i].Name < locations[j].Name })

	return ""
}

func buildSortedLocationList(timetables []wttxml.Timetable) []string {

	locations := map[string]bool{}

	for _, timetable := range timetables {
		for _, trip := range timetable.Trips.Trip {
			locations[trip.Location] = true
		}

	}

	uniqueLocations := make([]string, 0, len(locations))

	for name, _ := range locations {
		uniqueLocations = append(uniqueLocations, name)
	}

	sort.Strings(uniqueLocations)

	return uniqueLocations
}

//Data Structures for Stop List
type LocationStopList struct {
	Stops []LocationStop
}

func (a LocationStopList) Len() int           { return len(a.Stops) }
func (a LocationStopList) Swap(i, j int)      { a.Stops[i], a.Stops[j] = a.Stops[j], a.Stops[i] }
func (a LocationStopList) Less(i, j int) bool { return a.Stops[i].Time() < a.Stops[j].Time() }

type LocationStop struct {
	Headcode       string
	Origin         string
	Destination    string
	Arrival        *uint
	Departure      *uint
	Platform       string
	Updated        bool
	Arrived        bool
	Departed       bool
	ActualPlatform string
	DelaySeconds   int
}

func (a LocationStop) Time() uint {
	if a.Departure != nil {
		return *a.Departure
	}
	return *a.Arrival
}

func (a LocationStop) FormatArrival() string {
	return formatOptionalTime(a.Arrival)
}

func (a LocationStop) FormatDeparture() string {
	return formatOptionalTime(a.Departure)
}

const cancelledAfter uint = 3
const hiddenAfter int = 10

func (a LocationStop) OnTimeMessage() string {
	if !a.Updated && currentClock.Clock > (a.Time()+(cancelledAfter*60)) {
		return "Cancelled"
	} else if a.Departed {
		return "Departed"
	} else if a.Arrived {
		return "Arrived"
	} else if (a.DelaySeconds / 60) == 0 {
		return "On Time"
	} else if a.DelaySeconds > 0 {
		return fmt.Sprintf("Delayed by %d minutes", a.DelaySeconds/60)
	} else {
		return fmt.Sprintf("Early by %d minutes", (a.DelaySeconds*-1)/60)
	}

}

func (a LocationStop) HideAfter() int {

	/*if *showAll {
		return math.MaxInt32
	}*/

	if a.Departed {
		return 0
	} else {
		result := int(a.Time()) + (hiddenAfter * 60)
		if a.DelaySeconds > 0 {
			result += a.DelaySeconds
		}
		return result
	}

	return math.MaxInt32
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func isTIPLOC(id string) bool {
	return (len(id) <= 7) && (strings.ToUpper(id) == id)
}

//Build the list of stops for every location
func buildLocationStopList(locations []string, timetables []wttxml.Timetable) map[string]*LocationStopList {

	locationStops := map[string]*LocationStopList{}

	for _, loc := range locations {
		locationStops[loc] = &LocationStopList{}
	}

	tiplocs := readTIPLOCs()

	for _, timetable := range timetables {

		for i, _ := range timetable.Trips.Trip {
			trip := timetable.Trips.Trip[i]

			if !wttxml.SBool(trip.IsPassTime) && !(wttxml.SBool(trip.SetDownOnly) && trip.ArrTime == 0) && !(trip.DepPassTime == 0 && trip.ArrTime == 0) {
				stop := LocationStop{timetable.ID, timetable.OriginName, timetable.DestinationName, nil, nil, trip.Platform, false, false, false, "", 0}

				if isTIPLOC(stop.Origin) {
					stop.Origin = strings.Title(strings.ToLower(tiplocs[stop.Origin]))
				}
				if isTIPLOC(stop.Destination) {
					stop.Destination = strings.Title(strings.ToLower(tiplocs[stop.Destination]))
				}

				if stop.Destination == "" {
					stop.Origin = ""
					stop.Destination = timetable.Description
				}

				if trip.ArrTime != 0 {
					stop.Arrival = &trip.ArrTime
				}

				if !wttxml.SBool(trip.SetDownOnly) {
					if trip.ArrTime != 0 && trip.DepPassTime == 0 {
						stop.Departure = stop.Arrival
					} else {
						stop.Departure = &trip.DepPassTime
					}
				}

				stopList := locationStops[trip.Location]
				stopList.Stops = append(stopList.Stops, stop)
				locationStops[trip.Location] = stopList

			}

		}

	}

	return locationStops
}

//TIPLOC
func readTIPLOCs() map[string]string {
	tiplocs := make(map[string]string)

	f, err := os.Open(localTemplatePath("tiploc.bin"))
	if err != nil {
		//println(err.Error())
		return tiplocs
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	dec := msgpack.NewDecoder(buf)

	err = dec.Decode(&tiplocs)
	if err != nil {
		//println(err.Error())
	}
	return tiplocs
}
