package main

import (
	"fmt"
	"net/http"
)

func hello_response(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world, this machine has become now a server!")
}

func wkb_response(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Now we need a sample data!")
}

func wkt_response(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Now we need wkt as well")
}

func metadata_response(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Now we need metadata as well")
}

func main() {
	http.HandleFunc("/wkb", wkb_response)
	http.HandleFunc("/wkt", wkt_response)
	http.HandleFunc("/metadata", metadata_response)
	http.HandleFunc("/", hello_response)
	http.ListenAndServe(":8000", nil)
}
