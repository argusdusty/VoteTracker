package ExampleSource

import (
	"github.com/argusdusty/VoteTracker/Sources"
)

var (
	Updater = Sources.FileUpdater{
		Source: &Sources.ExampleSource{
			Param: "string",
		},
	}
)
