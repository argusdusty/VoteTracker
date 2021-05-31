package Sources

import (
	"fmt"
	"os"

	"github.com/argusdusty/VoteTracker/Utils"
)

type Source interface {
	Load(Utils.Summary) (Utils.Summary, error)
}

type FileUpdater struct {
	Source Source
}

func (U FileUpdater) Update(dst string) (result Utils.Summary, err error) {
	var Data Utils.Summary
	_, err = Utils.LoadSummary(dst, &Data)
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
