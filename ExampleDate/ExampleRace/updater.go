package ExampleRace

import (
	"VoteTracker/ExampleDate/ExampleRace/ExampleSource"
	. "VoteTracker/Utils"
)

var (
	Updaters = CombinedUpdater([]UpdaterDst{
		{ExampleSource.Updater, "ExampleSource"},
	})
)
