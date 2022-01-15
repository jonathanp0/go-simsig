package wttxml

import "encoding/xml"

// SavedTimetable.xml root element
type SimSigTimetable struct {
	XMLName    xml.Name `xml:"SimSigTimetable"`
	Timetables Timetables
}

type Timetables struct {
	Timetable []Timetable
}

type Timetable struct {
	XMLName         xml.Name `xml:"Timetable"`
	ID              string
	UID             string
	OriginName      string
	DestinationName string
	Description     string
	Trips           Trips
}

type Trips struct {
	XMLName xml.Name `xml:"Trips"`
	Trip    []Trip
}

type Trip struct {
	XMLName     xml.Name `xml:"Trip"`
	Location    string
	IsPassTime  string //-1 / True
	SetDownOnly string //-1 / True
	DepPassTime uint
	ArrTime     uint
	Platform    string
}

//In some timetable files "-1" is used instead of "True"...
func SBool(s string) bool {
	if s == "-1" || s == "True" {
		return true
	}
	return false
}
