package config

import (
	"log"
	"os"
	"strings"

	"github.com/naoina/toml"
)

type Configuration struct {
	BlobServerURL string

	ConfigDir      string
	CAKeyFile      string
	CACertFile     string
	StorageKeyFile string

	IndexDir string
}

func ReadConfig(path string) *Configuration {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	defer f.Close()

	c := &Configuration{}
	err = toml.NewDecoder(f).Decode(c)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	c.ConfigDir = strings.Replace(c.ConfigDir, "~", os.Getenv("HOME"), 1)

	return c
}
