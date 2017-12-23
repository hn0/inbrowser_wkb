package main

import (
	"db"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/lukeroth/gdal"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type app struct {
	database    *db.DB
	gdalsrs_bin string
}

var application *app

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	binary.Write(w, binary.LittleEndian, data)

	fmt.Printf("Binary request took %d ms (%d records delivered; content-length: %d)\n", time.Now().Sub(now)/100000, cnt, len(data))
}

func wkt_response(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	cnt, model := application.database.GetGeometry()
	var resp []interface{}
	if cnt == 0 {
		resp = make([]interface{}, 1)
		resp[0] = map[string]string{"error": "No records found!"}
	} else {
		resp = make([]interface{}, cnt)
		for i, _ := range resp {
			r := model.GetRecord(i)
			// read the geometry!
			var wkt string = "POINT EMPTY"

			ref := new(gdal.SpatialReference)
			b := r["GEOMETRY"].([]uint8)
			if geo, err := gdal.CreateFromWKB(b, *ref, len(b)); err == nil {
				wkt, _ = geo.ToWKT()
			}

			resp[i] = struct {
				Id  string
				WKT string
			}{
				strconv.FormatInt((*(r["ogc_fid"].(*interface{}))).(int64), 12),
				wkt,
			}
		}
	}
	close_request_json(resp, w)
	tmp, _ := json.Marshal(resp)
	fmt.Printf("Text request took %d ms (%d records delivered; content-length: %d)\n", time.Now().Sub(now)/100000, cnt, len(tmp))
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

func info_response(w http.ResponseWriter, r *http.Request) {
	var resp []interface{}
	resp = make([]interface{}, 1)
	resp[0] = struct {
		Source  string
		Geomcnt int64
		EPSG    int
		Proj    string
	}{
		application.database.GetSource(),
		application.database.GetCount(),
		application.GetEPSGCode(),
		application.GetDSProjection(),
	}
	close_request_json(resp, w)
}

func close_request_json(values []interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(values)
}

// helper functions

// returns proj4js datasource projection
func (a *app) GetDSProjection() string {
	opts := []string{"-o", "proj4", a.database.GetSource()}
	if out, err := exec.Command(a.gdalsrs_bin, opts...).Output(); err == nil {
		return strings.TrimSpace(string(out))
	} else {
		fmt.Println(err)
	}
	return "-"
}

// gives its best to extract epsg code
func (a *app) GetEPSGCode() int {

	opts := []string{"-o", "epsg", a.database.GetSource()}
	if out, err := exec.Command(a.gdalsrs_bin, opts...).Output(); err == nil {
		if string(out) == "EPSG:-1\n" {
			opts[1] = "wkt"
			if out, err = exec.Command(a.gdalsrs_bin, opts...).Output(); err != nil {
				return -1
			}
		}
		var parsval string
		cpos := strings.LastIndex(string(out), "EPSG") + 4
		// fmt.Println(string(string(out)[cpos]))
		switch c := string(out)[cpos]; string(c) {
		case ":":
			parsval = string(out)[cpos+1:]
		case "\"":
			// THERE IS AN ERROR HERE, FOR EG 900913, but is irrelevant for this purpose
			parsval = string(out)[cpos+3 : cpos+7]
		default:
			parsval = "-1"
		}
		if retval, err := strconv.Atoi(parsval); err == nil {
			return retval
		}
	} else {
		fmt.Println(err)
	}

	return -1
}

func main() {

	if len(os.Args) != 3 {
		log.Fatal("Usage: app sample_database_path static_pages_dir")
		log.Fatal("where: sample_database_path is file pointer to sqlite database holding sample geometry")
		log.Fatal("and static_pages_dir is path to the folder containing html/javascript files")
		log.Fatal("Exiting...")
	}

	application = new(app)
	application.gdalsrs_bin = "/usr/bin/gdalsrsinfo"
	if application.database = db.GetConn(os.Args[1]); application.database == nil {
		log.Fatal("Could not connect to sample database! Exiting ...")
	}

	http.HandleFunc("/wkb", wkb_response)
	http.HandleFunc("/wasm", wkb_response)
	http.HandleFunc("/wkt", wkt_response)
	http.HandleFunc("/metadata", metadata_response)
	http.HandleFunc("/geo", info_response)
	http.Handle("/", http.FileServer(http.Dir(os.Args[2])))
	http.ListenAndServe(":8000", nil)
	fmt.Println("up and running")
}
