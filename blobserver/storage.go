package main

import (
	"io"
	"log"
	"os"
	"path"
)

type DataStore struct {
	Path string
}

func OpenDataStore(path string) *DataStore {
	os.MkdirAll(path, 0700)
	return &DataStore{
		Path: path,
	}
}

func (s *DataStore) Put(key string, content io.Reader) error {
	p := s.pathForKey(key)
	os.MkdirAll(path.Dir(p), 0700)

	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	defer f.Close()

	_, err = io.Copy(f, content)
	log.Println("New blob", key)
	return err
}

func (s *DataStore) pathForKey(key string) string {
	p1 := key[:2]
	p2 := key[2:4]
	return path.Join(s.Path, p1, p2, key)
}

func (s *DataStore) Exists(key string) bool {
	_, err := os.Stat(s.pathForKey(key))
	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		log.Fatalln(err)
	}
	return false
}

func (s *DataStore) Get(key string, w io.Writer) error {
	f, err := os.Open(s.pathForKey(key))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}
