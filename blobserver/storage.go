package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

type DataStore struct {
	blobsPath string
	tempDir   string
	stateHash string
}

func OpenDataStore(storagePath string) *DataStore {
	bp := path.Join(storagePath, "blobs")
	tmp := path.Join(storagePath, "tmp")
	os.MkdirAll(bp, 0700)
	os.MkdirAll(tmp, 0700)
	s := &DataStore{
		blobsPath: bp,
		tempDir:   tmp,
	}
	s.saveStateHash()
	return s
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

func (s *DataStore) Get(key string) (io.ReadCloser, error) {
	return os.Open(s.pathForKey(key))
}

func (s *DataStore) walkKeys(keyHandler func(key string)) error {
	return filepath.Walk(s.blobsPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && len(info.Name()) > 4 {
			keyHandler(info.Name())
		}
		return nil
	})
}

func (s *DataStore) saveStateHash() string {
	hash := md5.New()
	s.walkKeys(func(key string) {
		hash.Write([]byte(key))
	})

	s.stateHash = hex.EncodeToString(hash.Sum(nil))
	return s.stateHash
}

func (s *DataStore) getStateHash() string {
	return s.stateHash
}
