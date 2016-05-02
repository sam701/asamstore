package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	http.HandleFunc("/hello", hello)

	configDir := path.Join(os.Getenv("HOME"), ".config/asamstore")

	caCert, err := ioutil.ReadFile(path.Join(configDir, "asamstore.cert.pem"))
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(path.Join(configDir, "server.cert.pem"), path.Join(configDir, "server.priv.pem"))
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientCAs:          caCertPool,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      "127.0.0.1:9000",
		TLSConfig: tlsConfig,
	}

	err = server.ListenAndServeTLS(
		path.Join(configDir, "server.cert.pem"),
		path.Join(configDir, "server.priv.pem"))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there!")
}
