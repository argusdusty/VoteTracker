package Utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"
)

type Races []Race

type Race struct {
	Name string `json:"name"`
	Race string `json:"race"`
}

func (R Races) GetText(date string, other bool) []string {
	t := make([]string, len(R)+1)
	if other {
		t[0] = "Other Races:"
	} else {
		t[0] = "Races:"
	}
	for i, r := range R {
		t[i+1] = fmt.Sprintf("%s: <a href=\"/%s/%s\">%s</a>", r.Name, date, r.Race, r.Race)
	}
	return t
}

func LoadRaces(dst string, R *Races) (error, time.Time) {
	f, err := os.Open(path.Join(dst, "races.json"))
	if err != nil {
		return err, time.Time{}
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return json.NewDecoder(f).Decode(R), time.Time{}
	}
	return json.NewDecoder(f).Decode(R), stat.ModTime()
}

func loadRaces(vars map[string]string) (interface{}, time.Time, error) {
	dst := path.Join(vars["date"])
	var R Races
	err, modtime := LoadRaces(dst, &R)
	fmt.Println(err, modtime, R)
	return R, modtime, err
}

func writeRaces(w http.ResponseWriter, v interface{}, vars map[string]string) {
	R := v.(Races)
	WriteHtmlHeader(w, "Races: "+ConvertDate(vars["date"]), false, false, false)
	WriteHtmlLines(w, R.GetText(vars["date"], false))
	WriteHtmlFooter(w)
}

func RacesHandler(w http.ResponseWriter, r *http.Request) {
	ValueHandler(w, r, "Races", loadRaces, writeRaces)
}
