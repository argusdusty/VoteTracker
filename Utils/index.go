package Utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Index []DateRace

type DateRace struct {
	Name string
	Date string
	Race string
}

func (I Index) GetText() []string {
	t := make([]string, len(I)+1)
	t[0] = "Featured Races:"
	for i, r := range I {
		t[i+1] = fmt.Sprintf("<a href=\"/%s/%s\">%s</a>", r.Date, r.Race, r.Name)
	}
	return t
}

func LoadIndex(I *Index) (modtime time.Time, err error) {
	f, err := os.Open("index.json")
	if err != nil {
		return time.Time{}, err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return time.Time{}, json.NewDecoder(f).Decode(I)
	}
	return stat.ModTime(), json.NewDecoder(f).Decode(I)
}

func loadIndex(vars map[string]string) (interface{}, time.Time, error) {
	var I Index
	modtime, err := LoadIndex(&I)
	return I, modtime, err
}

func writeIndex(w http.ResponseWriter, v interface{}, vars map[string]string) {
	I := v.(Index)
	WriteHtmlHeader(w, "Featured Races:", false, false, false)
	WriteHtmlLines(w, I.GetText())
	WriteHtmlFooter(w)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	ValueHandler(w, r, "Index", loadIndex, writeIndex)
}
