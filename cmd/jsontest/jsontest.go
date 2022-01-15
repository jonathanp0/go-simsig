package main

import (
	"encoding/json"
	"fmt"

	"github.com/jonathanp0/go-simsig/gateway"
)

//Attempt to decodea train movement message
func main() {

	b := []byte(`{"train_location":{"headcode":"2C29","uid":"4","action":"pass","location":"S717","platform":"","time":50893,"aspPass":6,"aspAppr":6}}`)
	var m gateway.TrainMovementMessage
	err := json.Unmarshal(b, &m)

	if err != nil {
		panic(err)
	}
	fmt.Println("Decoded Message for ", m.TrainLocation.Headcode)
}
