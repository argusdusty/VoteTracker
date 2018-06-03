package main

import (
	"VoteTracker/ExampleDate/ExampleRace"
	. "VoteTracker/Utils"
	"fmt"
	"time"
)

var (
	Updaters = MultiUpdater([]UpdaterDst{
		UpdaterDst{ExampleRace.Updaters, "ExampleRace"},
	})
	Frequency = 10 * time.Second
)

func doupdate() {
	err := Updaters.Update(".")
	if err != nil {
		fmt.Println("Update error:", err)
	}
}

func main() {
	doupdate()
	for range time.Tick(Frequency) {
		doupdate()
	}
}
