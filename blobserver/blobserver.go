package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var store *DataStore

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config := readConfig(os.Args[1])
	store = OpenDataStore(config.StorageDir)

	http.HandleFunc("/blob/", handleBlob)
	startHttpsServer(config)
}

func handleBlob(w http.ResponseWriter, r *http.Request) {
	up := strings.Split(r.URL.Path, "/")
	if len(up) < 3 {
		w.WriteHeader(400)
		io.WriteString(w, "No key provided")
		return
	}
	key := up[2]
	if len(key) < 10 {
		w.WriteHeader(400)
		io.WriteString(w, "Bad key: "+key+" in "+r.URL.Path)
		return
	}
	switch r.Method {
	case "GET":
		if r.URL.Query().Get("ifExists") != "" {
			if store.Exists(key) {
				w.WriteHeader(204)
			} else {
				w.WriteHeader(404)
			}
		} else {
			err := store.Get(key, w)
			if err != nil {
				if os.IsNotExist(err) {
					w.WriteHeader(404)
				} else {
					log.Println("Cannot write key", key, err)
					w.WriteHeader(500)
				}
			}
		}
	case "POST":
		err := store.Put(key, r.Body)
		if err != nil {
			log.Println("Cannot post key", key, err)
			w.WriteHeader(500)
		} else {
			w.WriteHeader(204)
		}
	default:
		w.WriteHeader(404)
		io.WriteString(w, "Unknown method "+r.Method)
	}
}
