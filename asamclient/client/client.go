package client

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/sam701/asamstore/asamclient/config"
)

type BlobStorageClient struct {
	url    string
	client *http.Client
}

func NewClient(c *config.Configuration) *BlobStorageClient {
	configDir := c.ConfigDir

	certFile := path.Join(configDir, c.CACertFile)
	cert, err := tls.LoadX509KeyPair(certFile, path.Join(configDir, c.CAKeyFile))
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := ioutil.ReadFile(certFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return &BlobStorageClient{
		url:    strings.TrimRight(c.BlobServerURL, "/") + "/blob/",
		client: &http.Client{Transport: transport},
	}
}

func (c *BlobStorageClient) Put(key string, content io.Reader) {
	resp, err := c.client.Get(c.url + key + "?ifExists=true")
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	switch resp.StatusCode {
	case 404:
		// not exists

		// send content to the server
		resp, err = c.client.Post(c.url+key, "application/octet-stream", content)
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		if resp.StatusCode != 204 {
			handleUnexpectedResponse(resp)
		}
	case 204:
		// content already exists
		log.Println("Key", key, "already exists")
	default:
		handleUnexpectedResponse(resp)
	}
}

func handleUnexpectedResponse(resp *http.Response) {
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	log.Fatalf("Status %d, message: %s\n", resp.StatusCode, string(bb))
}

func (c *BlobStorageClient) Get(key string, w io.Writer) bool {
	resp, err := c.client.Get(c.url + key)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	if resp.StatusCode == 404 {
		return false
	}
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	return true
}
