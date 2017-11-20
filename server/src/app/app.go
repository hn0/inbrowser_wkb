package main

import (
	"db"
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
	fmt.Fprintln(w, "Now we need metadata as well")
	cnt, model := application.database.GetMetadata()
	if cnt != 0 {
		cnt = 1

	}

	// close_request([]interface{"error"}, w)
	// var resp []interface{}
	// resp = make([]interface{}, cnt)
	//   should go something like this!?
}

func close_request(values []interface{}, w http.ResponseWriter) {
	// TODO: return a simple json response?!
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
