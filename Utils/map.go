package Utils

import (
	"net/http"
	"path"
	"time"
)

func loadMap(vars map[string]string) (interface{}, time.Time, error) {
	var S Summary
	dst := path.Join(vars["date"], vars["race"])
	if vars["source"] != "" {
		dst = path.Join(dst, vars["source"])
	}
	err, modtime := LoadSummary(dst, &S)
	if S.Regions == nil || len(S.Regions) == 0 {
		return nil, modtime, err
	}
	return S, modtime, err
}

func writeMap(w http.ResponseWriter, v interface{}, vars map[string]string) {
	S := v.(Summary)
	WriteHtmlHeader(w, S.Name+" Map", S.PortionComplete != 1, false, true)
	WriteMapScript(w)
	WriteHtmlFooter(w)
}

func MapHandler(w http.ResponseWriter, r *http.Request) {
	ValueHandler(w, r, "Map", loadMap, writeMap)
}
