package ExampleSource

import (
	. "VoteTracker/Sources"
)

var (
	Updater = FileUpdater{
		&ExampleSource{
			Param: "string",
		},
	}
)
