package put

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/codegangsta/cli"
)

func PutAction(c *cli.Context) error {
	configDir := path.Join(os.Getenv("HOME"), ".config/asamstore")

	cert, err := tls.LoadX509KeyPair(path.Join(configDir, "asamstore.cert.pem"), path.Join(configDir, "asamstore.priv.pem"))
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := ioutil.ReadFile(path.Join(configDir, "asamstore.cert.pem"))
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
	client := &http.Client{Transport: transport}

	// Do GET something
	resp, err := client.Get("https://localhost:9000/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))

	return nil
}
