package main

import (
	"db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type app struct {
	database *db.DB
}

var application *app

func hello_response(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world, this machine has become now a server!")
	fmt.Fprintln(w, "Serving sample data: "+application.database.GetSource())
}

func wkb_response(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Now we need a sample data!")
}

func wkt_response(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Now we need wkt as well")
}

func metadata_response(w http.ResponseWriter, r *http.Request) {
	cnt, model := application.database.GetMetadata()
	var resp []interface{}
	if cnt == 0 {
		resp = make([]interface{}, 1)
		resp[0] = map[string]string{"error": "No records found!"}
	} else {
		resp = make([]interface{}, cnt)
		for i, _ := range resp {
			resp[i] = model.GetRecord(i)
		}
	}
	close_request(resp, w)
}

func close_request(values []interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	json.NewEncoder(w).Encode(values)
}

func main() {

	if len(os.Args) != 2 {
		log.Fatal("A full path to sample database needs to be provided! Exiting ...")
	}

	application = new(app)
	if application.database = db.GetConn(os.Args[1]); application.database == nil {
		log.Fatal("Could not connect to sample database! Exiting ...")
	}

	http.HandleFunc("/wkb", wkb_response)
	http.HandleFunc("/wkt", wkt_response)
	http.HandleFunc("/metadata", metadata_response)
	http.HandleFunc("/", hello_response)
	http.ListenAndServe(":8000", nil)
	fmt.Println("up and running")
}
