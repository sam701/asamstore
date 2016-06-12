package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang/snappy"
)

func setupHttpHandlers() {
	http.HandleFunc("/blobs/", handleBlob)
	http.HandleFunc("/blobs/keys", getSortedKeys)
	http.HandleFunc("/blobs/keys/hash", writeStateHash)
	http.HandleFunc("/updateState", updateState)
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
			r, err := store.Get(key)
			if err != nil {
				if os.IsNotExist(err) {
					w.WriteHeader(404)
				} else {
					log.Println("Cannot write key", key, err)
					w.WriteHeader(500)
				}
			}
			defer r.Close()
			_, err = io.Copy(w, r)
			if err != nil {
				log.Println("ERROR: could not copy blob content", err)
			}
		}
	case "POST":
		err := store.Put(key, r.Body)
		if err != nil {
			log.Println("Cannot post key", key, err)
			w.WriteHeader(500)
			fmt.Fprintln(w, "Cannot post key", key, err)
		} else {
			w.WriteHeader(204)
		}
	default:
		w.WriteHeader(404)
		io.WriteString(w, "Unknown method "+r.Method)
	}
}

func getSortedKeys(w http.ResponseWriter, r *http.Request) {
	compressing := snappy.NewWriter(w)
	defer compressing.Close()

	store.walkKeys(func(key string) {
		fmt.Fprintln(compressing, key)
	})
}

func writeStateHash(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, store.getStateHash())
}

func updateState(w http.ResponseWriter, r *http.Request) {
	store.saveStateHash()
	go syncWithAllRemotes()
	w.WriteHeader(204)
}
