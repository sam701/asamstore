package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

func startHttpsServer(config *configuration) {
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
		Addr:      config.ServerAddress,
		TLSConfig: tlsConfig,
	}

	log.Println("Listening on", config.ServerAddress)
	err = server.ListenAndServeTLS(serverCert, serverKey)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
