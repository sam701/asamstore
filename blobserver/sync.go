package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/golang/snappy"
)

var blobServerSyncMutex sync.Mutex

func syncWithAllRemotes() {
	blobServerSyncMutex.Lock()
	defer blobServerSyncMutex.Unlock()

	for name, baseUrl := range config.Remotes {
		s := &remoteSync{name, baseUrl}
		s.sync()
	}
}

type remoteSync struct {
	name string
	url  string
}

func (s *remoteSync) retrieveRemoteBlob(key string) error {
	res, err := tlsHttpCilent.Get(s.url + "/blobs/" + key)
	if err != nil {
		return err
	}
	err = store.Put(key, res.Body)
	log.Println("Retrieved remote blob", key)
	return err
}

func (s *remoteSync) sendLocalBlob(key string) error {
	r, err := store.Get(key)
	if err != nil {
		return err
	}
	defer r.Close()
	res, err := tlsHttpCilent.Post(s.url+"/blobs/"+key, "application/octet-stream", r)
	if err != nil {
		return err
	}
	if res.StatusCode != 204 {
		return errors.New("Unexpected return code: " + res.Status)
	}
	log.Println("Sent local blob", key)
	return nil
}

func (s *remoteSync) needSync() bool {
	log.Println("Syncing with", s.name)
	res, err := tlsHttpCilent.Get(s.url + "/blobs/keys/hash")
	if err != nil {
		log.Println("Remote", s.name, "is not available")
		return false
	}

	defer res.Body.Close()
	bb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Cannot read blobs hash", err)
		return false
	}

	if string(bb) == store.getStateHash() {
		log.Println("Hashes are equal:", store.getStateHash())
		return false
	}
	return true
}

func (s *remoteSync) sync() {
	if !s.needSync() {
		return
	}

	log.Println("Syncing with", s.name)

	res, err := tlsHttpCilent.Get(s.url + "/blobs/keys")
	if err != nil {
		log.Println("Cannot get blobs refs from remote", s.name)
		return
	}

	decompressing := snappy.NewReader(res.Body)
	remoteKeyScanner := bufio.NewScanner(decompressing)
	store.walkKeys(func(localKey string) {
		if remoteKeyScanner.Scan() {
			remoteKey := remoteKeyScanner.Text()
			switch strings.Compare(localKey, remoteKey) {
			case -1:
				s.sendLocalBlob(localKey)
			case 1:
				s.retrieveRemoteBlob(remoteKey)
			}
		} else {
			s.sendLocalBlob(localKey)
		}
	})
	for remoteKeyScanner.Scan() {
		s.retrieveRemoteBlob(remoteKeyScanner.Text())
	}

}
