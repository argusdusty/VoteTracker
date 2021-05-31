package main

import (
	"fmt"
	"time"

	"github.com/argusdusty/VoteTracker/ExampleDate/ExampleRace"
	"github.com/argusdusty/VoteTracker/Utils"
)

var (
	Updaters = Utils.MultiUpdater([]Utils.UpdaterDst{
		{U: ExampleRace.Updaters, Dst: "ExampleRace"},
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
