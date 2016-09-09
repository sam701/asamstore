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
	cert, err := tls.LoadX509KeyPair(c.CACertFile(), c.CAKeyFile())
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := ioutil.ReadFile(c.CACertFile())
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
		url:    strings.TrimRight(c.BlobServerURL, "/"),
		client: &http.Client{Transport: transport},
		enc:    newEncrypter(c.BlobKeyFile()),
	}
}

func (c *BlobStorageClient) blobUrl(key string) string {
	return c.url + "/blobs/" + key
}

func (c *BlobStorageClient) Exists(ref schema.BlobRef) bool {
	key := string(ref)
	resp, err := c.client.Get(c.blobUrl(key) + "?ifExists=true")
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return resp.StatusCode == 204
}

func (c *BlobStorageClient) Put(ref schema.BlobRef, content io.Reader) {
	if !c.Exists(ref) {
		resp, err := c.client.Post(c.blobUrl(string(ref)), "application/octet-stream", c.enc.encryptingReader(content))
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		if resp.StatusCode != 204 {
			handleUnexpectedResponse(resp)
		}
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

func (c *BlobStorageClient) Get(ref schema.BlobRef) []byte {
	key := string(ref)
	resp, err := c.client.Get(c.blobUrl(key))
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	if resp.StatusCode == 404 {
		return nil
	}

	var buf bytes.Buffer
	copyAndVerify(&buf, c.enc.decryptingReader(resp.Body), ref)
	return buf.Bytes()
}

func (c *BlobStorageClient) GetSchema(ref schema.BlobRef) *schema.Schema {
	content := c.Get(ref)
	if content == nil {
		return nil
	} else {
		var s schema.Schema
		err := json.NewDecoder(bytes.NewReader(content)).Decode(&s)
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		return &s
	}
}

func (c *BlobStorageClient) UpdateServerState() {
	resp, err := c.client.Get(c.url + "/updateState")
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	if resp.StatusCode != 204 {
		log.Fatalln("Unexpected response for /updateState", resp.StatusCode)
	}
}
