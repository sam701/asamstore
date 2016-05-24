package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/sam701/asamstore/asamclient/config"
	"github.com/sam701/asamstore/asamclient/schema"
)

type BlobStorageClient struct {
	url    string
	client *http.Client
	enc    *encrypter
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
		enc:    newEncrypter(path.Join(configDir, c.BlobKeyFile)),
	}
}

func (c *BlobStorageClient) Put(ref schema.BlobRef, content io.Reader) {
	key := string(ref)
	resp, err := c.client.Get(c.url + key + "?ifExists=true")
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	switch resp.StatusCode {
	case 404:
		// not exists

		// send content to the server
		resp, err = c.client.Post(c.url+key, "application/octet-stream", c.enc.encryptingReader(content))
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		if resp.StatusCode != 204 {
			handleUnexpectedResponse(resp)
		}
	case 204:
		// content already exists
	default:
		handleUnexpectedResponse(resp)
	}
}

func (c *BlobStorageClient) PutSchema(s *schema.Schema) schema.BlobRef {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(s)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	ref := schema.GetBlobRefBytes(buf.Bytes())
	c.Put(ref, bytes.NewReader(buf.Bytes()))
	return ref
}

func handleUnexpectedResponse(resp *http.Response) {
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	log.Fatalf("Status %d, message: %s\n", resp.StatusCode, string(bb))
}

func (c *BlobStorageClient) Get(ref schema.BlobRef, w io.Writer) bool {
	key := string(ref)
	resp, err := c.client.Get(c.url + key)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	if resp.StatusCode == 404 {
		return false
	}
	_, err = io.Copy(w, c.enc.decryptingReader(resp.Body))
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	return true
}

func (c *BlobStorageClient) GetSchema(ref schema.BlobRef) *schema.Schema {
	var buf bytes.Buffer
	ok := c.Get(ref, &buf)
	if ok {
		var s schema.Schema
		err := json.NewDecoder(bytes.NewReader(buf.Bytes())).Decode(&s)
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		return &s
	} else {
		return nil
	}
}
