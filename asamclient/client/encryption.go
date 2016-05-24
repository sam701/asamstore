package client

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"

	"github.com/golang/snappy"
)

type encrypter struct {
	block cipher.Block
}

func newEncrypter(pwdFile string) *encrypter {
	data, err := ioutil.ReadFile(pwdFile)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	key, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return &encrypter{
		block: block,
	}
}

func (e *encrypter) encryptingReader(r io.Reader) io.Reader {
	var iv [aes.BlockSize]byte
	_, err := rand.Read(iv[:])
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	stream := cipher.NewCTR(e.block, iv[:])

	var buf bytes.Buffer
	_, err = buf.Write(iv[:])
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	sw := &cipher.StreamWriter{S: stream, W: &buf}
	cw := snappy.NewBufferedWriter(sw)

	_, err = io.Copy(cw, r)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	err = cw.Close()
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return &buf
}

func (e *encrypter) decryptingReader(r io.Reader) io.Reader {
	iv := make([]byte, aes.BlockSize)
	n, err := r.Read(iv)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	if n != len(iv) {
		log.Fatalln("Cannot read full IV, read just bytes:", n)
	}
	stream := cipher.NewCTR(e.block, iv)
	dr := &cipher.StreamReader{stream, r}

	return snappy.NewReader(dr)
}
