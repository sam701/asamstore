package config

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/naoina/toml"
)

type Configuration struct {
	BlobServerURL string

	CertificateDir string
	IndexDir       string
}

func (c *Configuration) CAKeyFile() string {
	return path.Join(c.CertificateDir, "asamstore.priv.pem")
}

func (c *Configuration) CACertFile() string {
	return path.Join(c.CertificateDir, "asamstore.cert.pem")
}
func (c *Configuration) BlobKeyFile() string {
	return path.Join(c.CertificateDir, "blob.key")
}

func ReadConfig(configPath string) *Configuration {
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	defer f.Close()

	c := &Configuration{}
	err = toml.NewDecoder(f).Decode(c)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	c.CertificateDir = strings.Replace(c.CertificateDir, "~", os.Getenv("HOME"), 1)
	c.IndexDir = strings.Replace(c.IndexDir, "~", os.Getenv("HOME"), 1)

	return c
}
