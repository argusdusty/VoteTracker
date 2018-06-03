package Utils

import (
	"fmt"
	"os"
	"path"
)

type Updater interface {
	Update(dst string) (Summary, error)
}

type UpdaterDst struct {
	U   Updater
	Dst string
}

type CombinedUpdater []UpdaterDst

func (U CombinedUpdater) Update(dst string) (Summary, error) {
	var Data Summary
	err, _ := LoadSummary(dst, &Data)
	if err != nil && !os.IsNotExist(err) {
		return Summary{}, err
	}
	Summaries := make([]Summary, len(U))
	for i, u := range U {
		Summaries[i], err = u.U.Update(path.Join(dst, u.Dst))
		if err != nil {
			return Summary{}, err
		}
	}
	result := CombineSummaries(Data.SummaryDefaults, Summaries...)
	if !Data.Equal(result) {
		err := result.SaveToFile(dst)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

type MultiUpdater []UpdaterDst

func (U MultiUpdater) Update(dst string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic:", r)
		}
	}()
	for _, u := range U {
		if _, err := u.U.Update(path.Join(dst, u.Dst)); err != nil {
			return err
		}
	}
	return nil
}
