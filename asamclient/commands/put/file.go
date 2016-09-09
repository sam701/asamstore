package put

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/sam701/asamstore/asamclient/schema"
)

func putFileParts(f *os.File, st os.FileInfo) schema.BlobRef {
	var pos int64 = 0
	fileSize := st.Size()
	parts := []*schema.BytesPart{}

	buf := make([]byte, schema.MaxFilePartSize)

	for pos < fileSize {
		n, err := f.Read(buf)
		if err == io.EOF {
			if n == 0 {
				break
			} else {
				err = nil
			}
		}
		if err != nil {
			log.Fatalln("ERROR", err)
		}

		content := buf[:n]
		ref := schema.GetBlobRefBytes(content)
		bsClient.Put(ref, bytes.NewReader(content))

		parts = append(parts, &schema.BytesPart{
			Size:       uint64(n),
			Offset:     uint64(pos),
			ContentRef: ref,
		})

		pos += int64(n)
	}

	return bsClient.PutSchema(getFileSchema(st, parts))
}

func getFileSchema(fi os.FileInfo, parts []*schema.BytesPart) *schema.Schema {
	s := schema.NewSchema(schema.ContentTypeFile)
	s.FileName = fi.Name()
	s.UnixPermission = fmt.Sprintf("%#o", fi.Mode())
	s.UnixMtime = fi.ModTime().Format(time.RFC3339)
	s.FileParts = parts
	return s
}
