package main

import (
	"log"
	"net/http"
	"net/http/pprof"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	err := http.ListenAndServe(":80", mux)

	if err != nil {
		log.Fatal(err)
	}

}
