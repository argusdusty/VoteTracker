package main

import (
	. "VoteTracker/Utils"
	"crypto/tls"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"path"
)

func RunAutocertServer() {
	certManager := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("certs"),
	}

	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
	panic(server.ListenAndServeTLS("", ""))
}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Favicon handler:", r.RemoteAddr, r.UserAgent())
	http.ServeFile(w, r, "favicon.ico")
}

func TopoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Topo handler:", r.RemoteAddr, r.UserAgent())
	vars := mux.Vars(r)
	http.ServeFile(w, r, path.Join(vars["date"], vars["race"], "topo.json"))
}

func ForecastJsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ForecastJs handler:", r.RemoteAddr, r.UserAgent())
	http.ServeFile(w, r, "forecast.js")
}

func MapJsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MapJs handler:", r.RemoteAddr, r.UserAgent())
	http.ServeFile(w, r, "map.js")
}

func TypeHandler(w http.ResponseWriter, r *http.Request) {
	switch mux.Vars(r)["type"] {
	case "", "summary", "Summary":
		SummaryHandler(w, r)
	case "forecast", "Forecast":
		ForecastHandler(w, r)
	case "map", "Map":
		MapHandler(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true).Host("{subdomain}.{domain}.{tld}").Subrouter()
	router.PathPrefix("/{date}/{race}/sources/{source}/{type}.{format}").HandlerFunc(TypeHandler)
	router.PathPrefix("/{date}/{race}/Sources/{source}/{type}.{format}").HandlerFunc(TypeHandler)
	router.PathPrefix("/{date}/{race}/sources/{source}/{type}").HandlerFunc(TypeHandler)
	router.PathPrefix("/{date}/{race}/Sources/{source}/{type}").HandlerFunc(TypeHandler)
	router.PathPrefix("/{date}/{race}/sources/{source}.{format}").HandlerFunc(TypeHandler)
	router.PathPrefix("/{date}/{race}/Sources/{source}.{format}").HandlerFunc(TypeHandler)
	router.PathPrefix("/{date}/{race}/sources/{source}").HandlerFunc(SummaryHandler)
	router.PathPrefix("/{date}/{race}/Sources/{source}").HandlerFunc(SummaryHandler)
	router.PathPrefix("/{date}/{race}/topo.json").HandlerFunc(TopoHandler)
	router.PathPrefix("/{date}/{race}/{type}.{format}").HandlerFunc(TypeHandler)
	router.PathPrefix("/{date}/{race}/{type}").HandlerFunc(TypeHandler)
	router.PathPrefix("/{date}/{race}.{format}").HandlerFunc(TypeHandler)
	router.PathPrefix("/{date}/{race}").HandlerFunc(SummaryHandler)
	router.PathPrefix("/favicon.ico").HandlerFunc(FaviconHandler)
	router.PathPrefix("/forecast.js").HandlerFunc(ForecastJsHandler)
	router.PathPrefix("/map.js").HandlerFunc(MapJsHandler)
	router.PathPrefix("/{date}").HandlerFunc(RacesHandler)
	router.PathPrefix("/").HandlerFunc(IndexHandler)
	http.Handle("/", router)
	RunAutocertServer()
}
