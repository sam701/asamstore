package initialize

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"path"
	"time"

	"github.com/codegangsta/cli"
)

func Initialize(c *cli.Context) error {
	configDir := path.Join(os.Getenv("HOME"), ".config/asamstore")
	os.MkdirAll(configDir, 0700)

	caKey := createPrivateKey(path.Join(configDir, "asamstore.priv.pem"))
	caCert := createCertificate(caKey, nil, path.Join(configDir, "asamstore.cert.pem"))

	serverKey := createPrivateKey(path.Join(configDir, "server.priv.pem"))
	createCertificate(serverKey, caCert, path.Join(configDir, "server.cert.pem"))

	return nil
}

func createPrivateKey(saveToPath string) *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	savePem(saveToPath, x509.MarshalPKCS1PrivateKey(key), "RSA PRIVATE KEY")
	return key
}

func createCertificate(key *rsa.PrivateKey, parentCert *x509.Certificate, saveToPath string) *x509.Certificate {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}
	isCA := parentCert == nil
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"asamstore"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(20, 0, 0),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,

		IsCA: isCA,
	}

	if isCA {
		template.KeyUsage |= x509.KeyUsageCertSign
		parentCert = &template
	} else {
		template.ExtKeyUsage = append(template.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, parentCert, &key.PublicKey, key)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	savePem(saveToPath, derBytes, "CERTIFICATE")

	out, err := x509.ParseCertificate(derBytes)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return out
}

func savePem(pemFilePath string, derBytes []byte, keyType string) {
	f, err := os.OpenFile(pemFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	defer f.Close()
	err = pem.Encode(f, &pem.Block{Type: keyType, Bytes: derBytes})
	if err != nil {
		log.Fatalln("ERROR", err)
	}
}
