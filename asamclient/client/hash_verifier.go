package client

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"log"

	"github.com/sam701/asamstore/asamclient/schema"
)

type hashVerifier struct {
	hash hash.Hash
	dest io.Writer
}

func (v *hashVerifier) Write(p []byte) (int, error) {
	_, err := v.hash.Write(p)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return v.dest.Write(p)
}

func copyAndVerify(dest io.Writer, src io.Reader, ref schema.BlobRef) {
	v := &hashVerifier{sha1.New(), dest}
	_, err := io.Copy(v, src)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	sum := v.hash.Sum(nil)
	calculatedKey := hex.EncodeToString(sum)

	blobKey := string(ref)
	if calculatedKey != blobKey {
		log.Fatalf("Blob key %s was not equal the calculated one %s\n", blobKey, calculatedKey)
	}
}
