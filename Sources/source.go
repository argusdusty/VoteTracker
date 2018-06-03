package Sources

import (
	. "VoteTracker/Utils"
	"fmt"
	"os"
)

type Source interface {
	Load(Summary) (Summary, error)
}

type FileUpdater struct {
	Source Source
}

func (U FileUpdater) Update(dst string) (result Summary, err error) {
	var Data Summary
	err, _ = LoadSummary(dst, &Data)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	result, err = U.Source.Load(Data)
	if err != nil {
		return
	}
	result.SummaryDefaults = Data.SummaryDefaults
	if !Data.Equal(result) {
		err = result.SaveToFile(dst)
		if err != nil {
			fmt.Println(result, Data, dst, err)
			return
		}
	}
	return
}
