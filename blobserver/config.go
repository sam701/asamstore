package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type configuration struct {
	StorageDir    string `yaml:"storageDir"`
	ServerAddress string `yaml:"serverAddress"`

	Certificates certificateConfig `yaml:"certificates"`
	Remotes      map[string]string `yaml:"remotes"`
}

type certificateConfig struct {
	Path       string `yaml:"path"`
	CA         string `yaml:"ca"`
	ServerKey  string `yaml:"serverKey"`
	ServerCert string `yaml:"serverCert"`
}

func readConfig(pathToConfig string) *configuration {
	bb, err := ioutil.ReadFile(pathToConfig)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	var c configuration
	err = yaml.Unmarshal(bb, &c)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	c.Certificates.Path = strings.Replace(c.Certificates.Path, "~", os.Getenv("HOME"), 1)
	return &c
}
