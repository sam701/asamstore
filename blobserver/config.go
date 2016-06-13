package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

type configuration struct {
	StorageDir    string `yaml:"storageDir"`
	ServerAddress string `yaml:"serverAddress"`

	Remotes map[string]string `yaml:"remotes"`
}

func (c *configuration) certsPath() string {
	return path.Join(c.StorageDir, "certs")
}

func (c *configuration) CAPath() string {
	return path.Join(c.certsPath(), "ca.cert.pem")
}

func (c *configuration) ServerKeyPath() string {
	return path.Join(c.certsPath(), "server.priv.pem")
}

func (c *configuration) ServerCertPath() string {
	return path.Join(c.certsPath(), "server.cert.pem")
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
	c.StorageDir = strings.Replace(c.StorageDir, "~", os.Getenv("HOME"), 1)
	for k, v := range c.Remotes {
		c.Remotes[k] = strings.TrimRight(v, "/")
	}
	return &c
}
