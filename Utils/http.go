package Utils

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

var start_time time.Time

func init() {
	start_time = time.Now()
}

func ValueHandler(w http.ResponseWriter, r *http.Request, name string, loadValue func(map[string]string) (interface{}, time.Time, error), writeValue func(http.ResponseWriter, interface{}, map[string]string)) {
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			fmt.Println(r)
		}
	}()
	fmt.Println(name, "handler:", r.RemoteAddr, r.UserAgent())
	vars := mux.Vars(r)
	v, modtime, err := loadValue(vars)
	if os.IsNotExist(err) || v == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
		return
	}
	if err != nil {
		panic(err)
	}
	modt := r.Header.Get("If-Modified-Since")
	t, err := time.Parse(http.TimeFormat, modt)
	if err == nil && modtime.Before(t.Add(time.Second)) && start_time.UTC().Before(t.Add(time.Second)) {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
	format := vars["format"]
	if format == "" {
		r.ParseForm()
		formats := r.Form["format"]
		if len(formats) > 0 {
			format = formats[0]
		}
	}
	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		err = json.NewEncoder(w).Encode(v)
		if err != nil {
			panic(err)
		}
	default:
		writeValue(w, v, vars)
	}
}
