package Utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"
	"time"
)

type SummaryDefaults struct {
	Name               string            `json:"name"`
	PrefixText         []string          `json:"prefix_text"`
	SuffixText         []string          `json:"suffix_text"`
	VotePortionThresh  float64           `json:"vote_portion_thresh"`
	MapPortionThresh   float64           `json:"map_portion_thresh"`
	OddsThresh         float64           `json:"odds_thresh"`
	GraphPortionThresh float64           `json:"graph_portion_thresh"`
	GraphOddsThresh    float64           `json:"graph_odds_thresh"`
	Colors             map[string]string `json:"colors"`
	Priority           map[string]int    `json:"priority"`
	OtherRaces         Races             `json:"other_races,omitempty"`
	ShowDate           bool              `json:"show_date"`
	ShowMap            bool              `json:"show_map"`
}

type SummaryVote struct {
	Candidates      []SummaryCandidate `json:"candidates"`
	TotalVotes      uint64             `json:"total_votes"`
	PortionComplete float64            `json:"portion_complete"`
}

type SummaryCandidate struct {
	Candidate string `json:"candidate"`
	Votes     uint64 `json:"votes"`
}

type Summary struct {
	SummaryDefaults `json:",inline"`
	SummaryVote     `json:",inline"`
	Regions         map[string]SummaryRegion `json:"regions,omitempty"`
	Forecast        Forecast                 `json:"forecast,omitempty"`
	ForecastWeight  float64                  `json:"forecast_weight"`
}

type SummaryRegion struct {
	SummaryVote `json:",inline"`
	Name        string `json:"name"`
	Exclude     bool   `json:"exclude"`
}

type SummaryVotePriority struct {
	SummaryVote
	Priority map[string]int
}

func (S SummaryVotePriority) Len() int { return len(S.Candidates) }
func (S SummaryVotePriority) Less(i, j int) bool {
	return S.Candidates[i].Votes > S.Candidates[j].Votes || (S.Candidates[i].Votes == S.Candidates[j].Votes && S.Priority[S.Candidates[i].Candidate] < S.Priority[S.Candidates[j].Candidate])
}
func (S SummaryVotePriority) Swap(i, j int) {
	S.SummaryVote.Candidates[i], S.Candidates[j] = S.Candidates[j], S.Candidates[i]
}

func (S Summary) Sort() {
	sort.Stable(SummaryVotePriority{S.SummaryVote, S.Priority})
	for _, region := range S.Regions {
		sort.Stable(SummaryVotePriority{region.SummaryVote, S.Priority})
	}
	if S.Forecast != nil {
		S.Forecast.Sort()
	}
}

func (S Summary) GetText(date, race string) []string {
	var s []string
	s = append(s, S.PrefixText...)

	S.Sort()
	for _, c := range S.Candidates {
		_, ok := S.Priority[c.Candidate]
		if (S.TotalVotes != 0 && float64(c.Votes) >= S.VotePortionThresh*float64(S.TotalVotes)) || (S.TotalVotes == 0 && ok) {
			if S.TotalVotes == 0 {
				s = append(s, fmt.Sprintf("%s: %d (0.0%%)", c.Candidate, c.Votes))
			} else {
				s = append(s, fmt.Sprintf("%s: %d (%.1f%%)", c.Candidate, c.Votes, float64(c.Votes)/float64(S.TotalVotes)*100))
			}
		}
	}
	s = append(s, fmt.Sprintf("Total votes: %d", S.TotalVotes))
	if len(S.Candidates) >= 2 && S.TotalVotes != 0 {
		s = append(s, fmt.Sprintf("%s margin: %.2f%%", S.Candidates[0].Candidate, float64(S.Candidates[0].Votes-S.Candidates[1].Votes)/float64(S.TotalVotes)*100))
	}

	if S.PortionComplete > 0 {
		if S.PortionComplete == 1 {
			s = append(s, fmt.Sprint("100% complete"))
		} else {
			s = append(s, fmt.Sprintf("Estimated percent complete: %.2f%%", S.PortionComplete*100))
		}
	}

	s = append(s, "")
	if S.Forecast != nil {
		s = append(s, S.Forecast.GetText(S.OddsThresh)...)
		s = append(s, fmt.Sprintf("<a href=\"/%s/%s/forecast\">See forecast</a>", date, race))
	}

	s = append(s, "")
	if S.OtherRaces != nil {
		s = append(s, S.OtherRaces.GetText(date, true)...)
	}

	if S.ShowDate {
		s = append(s, fmt.Sprintf("<a href=\"/%s\">See all races</a>", date))
	}

	s = append(s, S.SuffixText...)
	return s
}

func (S Summary) SaveToFile(dst string) error {
	fmt.Println(dst, S.Candidates, S.TotalVotes)
	f, err := os.OpenFile(path.Join(dst, "summary.json"), os.O_WRONLY|os.O_CREATE, 0666)
	f.Truncate(0)
	f.Seek(0, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(S)
}

func (A Summary) Equal(B Summary) bool {
	if len(A.Candidates) != len(B.Candidates) {
		return false
	}
	for i, c := range A.Candidates {
		if B.Candidates[i] != c {
			return false
		}
	}
	if len(A.Forecast) != len(B.Forecast) {
		return false
	}
	for i, c := range A.Forecast {
		if B.Forecast[i].Candidate != c.Candidate && B.Forecast[i].ConcentrationParam != c.ConcentrationParam {
			return false
		}
	}
	return true
}

func LoadSummary(dst string, S *Summary) (error, time.Time) {
	f, err := os.Open(path.Join(dst, "summary.json"))
	if err != nil {
		return err, time.Time{}
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return json.NewDecoder(f).Decode(S), time.Time{}
	}
	return json.NewDecoder(f).Decode(S), stat.ModTime()
}

func loadSummary(vars map[string]string) (interface{}, time.Time, error) {
	var S Summary
	dst := path.Join(vars["date"], vars["race"])
	if vars["source"] != "" {
		dst = path.Join(dst, vars["source"])
	}
	err, modtime := LoadSummary(dst, &S)
	return S, modtime, err
}

func writeSummary(w http.ResponseWriter, v interface{}, vars map[string]string) {
	S := v.(Summary)
	WriteHtmlHeader(w, S.Name, S.PortionComplete != 1, false, S.ShowMap)
	WriteHtmlLines(w, S.GetText(vars["date"], vars["race"]))
	if S.ShowMap {
		WriteMapScript(w)
	}
	WriteHtmlFooter(w)
}

func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	ValueHandler(w, r, "Summary", loadSummary, writeSummary)
}

func CombineSummaries(Defaults SummaryDefaults, Ss ...Summary) Summary {
	var result Summary
	result.SummaryDefaults = Defaults

	// Combine SummaryVote
	var res = make(map[string]uint64)
	var totalTmp uint64
	var portionTotal uint64
	for _, s := range Ss {
		totalTmp = 0
		for _, c := range s.Candidates {
			totalTmp += c.Votes
			if c.Votes > res[c.Candidate] || res[c.Candidate] == 0 {
				res[c.Candidate] = c.Votes
			}
		}
		if totalTmp > result.TotalVotes {
			result.TotalVotes = totalTmp
		}
		if totalTmp >= portionTotal && s.PortionComplete > result.PortionComplete {
			result.PortionComplete = s.PortionComplete
			portionTotal = totalTmp
		}
	}

	// Combine Forecast
	for _, s := range Ss {
		if result.Forecast == nil {
			result.Forecast = s.Forecast
		} else if s.Forecast != nil && s.ForecastWeight > result.ForecastWeight {
			result.Forecast = s.Forecast
		}
	}

	// Combine Regions
	result.Regions = make(map[string]SummaryRegion)
	var rres = make(map[string]map[string]uint64)
	var rportionTotal = make(map[string]uint64)
	for _, s := range Ss {
		for k, region := range s.Regions {
			r, ok := result.Regions[k]
			if !ok {
				r.Name = region.Name
				r.Exclude = region.Exclude
				rres[k] = make(map[string]uint64)
			} else if !region.Exclude {
				r.Exclude = false
			}
			totalTmp = 0
			for _, c := range region.Candidates {
				totalTmp += c.Votes
				if c.Votes > rres[k][c.Candidate] || rres[k][c.Candidate] == 0 {
					rres[k][c.Candidate] = c.Votes
				}
			}
			if totalTmp > r.TotalVotes {
				r.TotalVotes = totalTmp
			}
			if totalTmp >= rportionTotal[k] && region.PortionComplete > r.PortionComplete {
				r.PortionComplete = region.PortionComplete
				rportionTotal[k] = totalTmp
			}
			if r.Name == "" {
				r.Name = region.Name
			}
			result.Regions[k] = r
		}
	}
	var totalCandidates = make(map[string]uint64)
	var total uint64
	for k, region := range result.Regions {
		totalTmp = 0
		for c, v := range rres[k] {
			region.Candidates = append(region.Candidates, SummaryCandidate{Candidate: c, Votes: v})
			totalTmp += v
			totalCandidates[c] += v
		}
		if totalTmp > region.TotalVotes {
			region.TotalVotes = totalTmp
		}
		total += region.TotalVotes
		result.Regions[k] = region
	}
	if total > result.TotalVotes {
		result.TotalVotes = total
	}

	totalTmp = 0
	for c, v := range res {
		if totalCandidates[c] > v {
			v = totalCandidates[c]
		}
		result.Candidates = append(result.Candidates, SummaryCandidate{Candidate: c, Votes: v})
		totalTmp += v
	}
	if totalTmp > result.TotalVotes {
		result.TotalVotes = totalTmp
	}
	result.Sort()

	return result
}
