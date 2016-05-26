package main

import (
	"io"
	"log"
	"os"
	"path"
)

type DataStore struct {
	blobsPath string
	tempDir   string
}

func OpenDataStore(storagePath string) *DataStore {
	bp := path.Join(storagePath, "blobs")
	tmp := path.Join(storagePath, "tmp")
	os.MkdirAll(bp, 0700)
	os.MkdirAll(tmp, 0700)
	return &DataStore{
		blobsPath: bp,
		tempDir:   tmp,
	}
}

func (s *DataStore) Put(key string, content io.Reader) error {
	tmpFilePath := path.Join(s.tempDir, key)

	err := s.writeTempFile(tmpFilePath, content)
	if err != nil {
		return err
	}

	p := s.pathForKey(key)
	os.MkdirAll(path.Dir(p), 0700)

	err = os.Rename(tmpFilePath, p)
	if err != nil {
		return err
	}

	log.Println("New blob", key)
	return nil
}

func (s *DataStore) writeTempFile(tmpFilePath string, content io.Reader) error {
	f, err := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	defer f.Close()

	_, err = io.Copy(f, content)
	return err
}

func (s *DataStore) pathForKey(key string) string {
	p1 := key[:2]
	p2 := key[2:4]
	return path.Join(s.blobsPath, p1, p2, key)
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
