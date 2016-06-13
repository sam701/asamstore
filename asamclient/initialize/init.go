package initialize

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path"
	"time"

	"github.com/codegangsta/cli"
	"github.com/sam701/asamstore/asamclient/config"
)

func Initialize(c *cli.Context) error {
	destinationDir := c.String("dest-dir")
	if destinationDir == "" {
		destinationDir = path.Join(os.Getenv("HOME"), ".config/asamstore")
	}
	if _, err := os.Stat(destinationDir); err == nil {
		fmt.Println("Directory", destinationDir, "already exists")
		cli.ShowCommandHelp(c, "init")
		return nil
	}
	os.MkdirAll(destinationDir, 0700)

	if c.Bool("client") {
		caKey := createPrivateKey(path.Join(destinationDir, "asamstore.priv.pem"))
		createCertificate(&caKey.PublicKey, nil, caKey, path.Join(destinationDir, "asamstore.cert.pem"))
		generateAndSaveAESKey(path.Join(destinationDir, "blob.key"))
	} else if c.Bool("blob-server") {
		cfg := config.ReadConfig(c.GlobalString("config"))
		caKey := readPrivateKey(cfg.CAKeyFile())
		caCert := readCertificate(cfg.CACertFile())
		serverKey := createPrivateKey(path.Join(destinationDir, "server.priv.pem"))
		createCertificate(&serverKey.PublicKey, caCert, caKey, path.Join(destinationDir, "server.cert.pem"))
		err := os.Link(cfg.CACertFile(), path.Join(destinationDir, "ca.cert.pem"))
		if err != nil {
			log.Fatalln("ERROR", err)
		}
	}

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

func readPrivateKey(privateKeyPath string) *rsa.PrivateKey {
	pemBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	block, _ := pem.Decode(pemBytes)
	derBytes := block.Bytes

	key, err := x509.ParsePKCS1PrivateKey(derBytes)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	return key
}

func createCertificate(
	signeePublicKey *rsa.PublicKey,
	parentCert *x509.Certificate,
	signerPrivateKey *rsa.PrivateKey,
	saveToPath string) *x509.Certificate {

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

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, parentCert, signeePublicKey, signerPrivateKey)
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

func readCertificate(pemFilePath string) *x509.Certificate {
	bb, err := ioutil.ReadFile(pemFilePath)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	block, _ := pem.Decode(bb)
	derBytes := block.Bytes

	out, err := x509.ParseCertificate(derBytes)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return out
}

func generateAndSaveAESKey(pathToSave string) {
	f, err := os.OpenFile(pathToSave, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}
	defer f.Close()

	keyBuf := make([]byte, 32)
	_, err = rand.Read(keyBuf)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	encodedKey := base64.StdEncoding.EncodeToString(keyBuf)
	_, err = io.WriteString(f, encodedKey)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
}
