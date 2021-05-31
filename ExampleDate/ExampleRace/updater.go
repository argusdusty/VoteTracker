package ExampleRace

import (
	"github.com/argusdusty/VoteTracker/ExampleDate/ExampleRace/ExampleSource"
	"github.com/argusdusty/VoteTracker/Utils"
)

var (
	Updaters = Utils.CombinedUpdater([]Utils.UpdaterDst{
		{U: ExampleSource.Updater, Dst: "ExampleSource"},
	})
)
