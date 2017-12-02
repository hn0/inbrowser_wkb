package main

import (
	"db"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type app struct {
	database *db.DB
}

var application *app

func hello_response(w http.ResponseWriter, r *http.Request) {
	var resp []interface{}
	resp = make([]interface{}, 1)
	resp[0] = map[string]string{
		"message":       "Hello world, this machine has become now a server!",
		"sample source": application.database.GetSource(),
	}
	close_request_json(resp, w)
}

func wkb_response(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	cnt, model := application.database.GetGeometry()
	var data []byte
	data = make([]byte, cnt*4+model.SizeOf())

	pos := 0
	for i := 0; i < cnt; i++ {
		r := model.GetRecord(i)
		binary.LittleEndian.PutUint32(data[pos:pos+4], uint32((*(r["ogc_fid"]).(*interface{})).(int64)))
		pos += 4
		geom := r["GEOMETRY"].([]byte)
		l := len(geom)
		binary.LittleEndian.PutUint32(data[pos:pos+4], uint32(l))
		pos += 4
		copy(data[pos:pos+l], geom)
		pos += l
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	binary.Write(w, binary.LittleEndian, data)

	fmt.Printf("Binary request took %d ms (%d records delivered; content-length: %d)\n", time.Now().Sub(now)/1000, cnt, len(data))
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
	close_request_json(resp, w)
}

func close_request_json(values []interface{}, w http.ResponseWriter) {
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
