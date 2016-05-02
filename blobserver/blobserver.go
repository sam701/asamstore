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
	config := readConfig(os.Args[1])

	http.HandleFunc("/hello", hello)

	configDir := config.Certificates.Path

	caCert, err := ioutil.ReadFile(path.Join(configDir, config.Certificates.CA))
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	serverCert := path.Join(configDir, config.Certificates.ServerCert)
	serverKey := path.Join(configDir, config.Certificates.ServerKey)
	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
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

	err = server.ListenAndServeTLS(serverCert, serverKey)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there!")
}
