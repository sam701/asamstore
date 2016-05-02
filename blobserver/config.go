package main

import (
	"log"
	"os"
	"strings"

	"github.com/naoina/toml"
)

type configuration struct {
	StorageDir string

	Certificates struct {
		Path       string
		CA         string
		ServerKey  string
		ServerCert string
	}
}

func readConfig(pathToConfig string) *configuration {
	f, err := os.Open(pathToConfig)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	defer f.Close()

	var c configuration
	err = toml.NewDecoder(f).Decode(&c)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	c.Certificates.Path = strings.Replace(c.Certificates.Path, "~", os.Getenv("HOME"), 1)
	return &c
}
