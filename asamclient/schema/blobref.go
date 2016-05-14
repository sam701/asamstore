package schema

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
)

type BlobRef string

func GetBlobRef(r io.Reader) BlobRef {
	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return BlobRef(hex.EncodeToString(h.Sum(nil)))
}

func GetBlobRefBytes(b []byte) BlobRef {
	return GetBlobRef(bytes.NewReader(b))
}
