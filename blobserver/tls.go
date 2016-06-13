package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

var tlsHttpCilent *http.Client

func startHttpsServer() {
	caCert, err := ioutil.ReadFile(config.CAPath())
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(config.ServerCertPath(), config.ServerKeyPath())
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
	err = server.ListenAndServeTLS(config.ServerCertPath(), config.ServerKeyPath())
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initTlsClient() {
	cert, err := tls.LoadX509KeyPair(config.ServerCertPath(), config.ServerKeyPath())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("client:", len(cert.Certificate))

	caCert, err := ioutil.ReadFile(config.CAPath())
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
	tlsHttpCilent = &http.Client{Transport: transport}
}
