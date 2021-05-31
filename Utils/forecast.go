package Utils

import (
	"fmt"
	"net/http"
	"path"
	"sort"
	"time"

	"gonum.org/v1/gonum/stat/distmv"
)

type Forecast []CandidateForecast

type CandidateForecast struct {
	Candidate          string  `json:"candidate"`
	ConcentrationParam float64 `json:"concentration_param"` // Dirichlet distribution params
	Odds               float64 `json:"odds"`
}

func (F Forecast) Sort() {
	sort.SliceStable(F, func(i, j int) bool { return F[i].ConcentrationParam > F[j].ConcentrationParam })
}

func (F Forecast) SetOdds(sims uint64) {
	hits := make([]uint64, len(F))
	sample := make([]float64, len(F))
	var mx float64
	var b int
	x := make([]float64, len(F))
	for i, c := range F {
		x[i] = c.ConcentrationParam
	}
	dist := distmv.NewDirichlet(x, nil)
	for i := uint64(0); i < sims; i++ {
		mx, b = 0, 0
		dist.Rand(sample)
		for j, r := range sample {
			if r > mx {
				mx = r
				b = j
			}
		}
		hits[b]++
	}
	for i, o := range hits {
		F[i].Odds = float64(o+1) / float64(sims+uint64(len(F)))
	}
}

func (F Forecast) Equal(F2 Forecast) bool {
	if len(F) != len(F2) {
		return false
	}
	for i, c := range F {
		if c.Candidate != F2[i].Candidate && c.ConcentrationParam != F2[i].ConcentrationParam {
			return false
		}
	}
	return true
}

func (F Forecast) GetText(thresh float64) []string {
	s := make([]string, 0, len(F))
	var t float64
	for _, c := range F {
		if c.Odds > thresh && c.Candidate != "" {
			s = append(s, fmt.Sprintf("%s win probability: %.1f%%", c.Candidate, c.Odds*100))
		}
		t += c.ConcentrationParam
	}
	if len(F) > 1 {
		s = append(s, fmt.Sprintf("%s forecast margin: %.1f%%", F[0].Candidate, (F[0].ConcentrationParam-F[1].ConcentrationParam)/t*100))
	}
	return s
}

func loadForecast(vars map[string]string) (interface{}, time.Time, error) {
	var S Summary
	dst := path.Join(vars["date"], vars["race"])
	if vars["source"] != "" {
		dst = path.Join(dst, vars["source"])
	}
	modtime, err := LoadSummary(dst, &S)
	if S.Forecast == nil {
		return nil, modtime, err
	}
	return S, modtime, err
}

func writeForecast(w http.ResponseWriter, v interface{}, vars map[string]string) {
	S := v.(Summary)
	WriteHtmlHeader(w, S.Name+" Forecast", S.PortionComplete != 1, true, false)
	WriteHtmlLines(w, S.Forecast.GetText(S.OddsThresh))
	WriteForecastScript(w)
	WriteHtmlFooter(w)
}

func ForecastHandler(w http.ResponseWriter, r *http.Request) {
	ValueHandler(w, r, "Forecast", loadForecast, writeForecast)
}
