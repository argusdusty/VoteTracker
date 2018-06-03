package Utils

import (
	"io"
	"net/http"
	"time"
)

type CacheType struct {
	Value    interface{}
	ModTime  time.Time
	LastLoad time.Time
}

var Cache = make(map[string]CacheType)

func LoadURL(url string, f func(io.Reader) (interface{}, error)) (val interface{}, err error) {
	var modTime time.Time
	v := Cache[url]
	if v.Value != nil {
		modTime = v.ModTime
	}
	val = v.Value
	if time.Now().Sub(v.LastLoad) < 2*time.Second { // Prevent flooding any service, reuse cache result with a couple seconds
		return
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	if !modTime.IsZero() {
		req.Header.Set("If-Modified-Since", modTime.Format(http.TimeFormat))
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 304 {
		val = v.Value
	} else if resp.StatusCode == 200 {
		var lastMod time.Time
		if len(resp.Header["Last-Modified"]) > 0 {
			lastMod, err = time.Parse(http.TimeFormat, resp.Header["Last-Modified"][0])
		}
		if lastMod.IsZero() && len(resp.Header["Date"]) > 0 {
			lastMod, err = time.Parse(http.TimeFormat, resp.Header["Date"][0])
		}
		if lastMod.IsZero() {
			lastMod = time.Now()
		}
		val, err = f(resp.Body)
		Cache[url] = CacheType{Value: val, ModTime: lastMod, LastLoad: time.Now()}
	}
	return
}
