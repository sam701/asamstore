package put

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/sam701/asamstore/asamclient/client"
	"github.com/sam701/asamstore/asamclient/config"
	"github.com/sam701/asamstore/asamclient/schema"
)

var bsClient *client.BlobStorageClient

func PutAction(c *cli.Context) error {
	cfg := config.ReadConfig(c.GlobalString("config"))
	bsClient = client.NewClient(cfg)

	contentPath := c.Args().First()
	if contentPath == "" {
		log.Fatalln("No ontent path provided")
	}
	putFile(contentPath)

	return nil
}

func putFile(filePath string) schema.BlobRef {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	if st.IsDir() {
		fis, err := f.Readdir(0)
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		entries := []schema.BlobRef{}
		for _, fi := range fis {
			ref := putFile(path.Join(filePath, fi.Name()))
			entries = append(entries, ref)
		}
		return bsClient.PutSchema(getDirSchema(st, entries))
	} else {
		contentRef := schema.GetBlobRef(f)
		f.Seek(0, 0)
		bsClient.Put(contentRef, f)

		return bsClient.PutSchema(getFileSchema(st, contentRef))
	}
}

func getDirSchema(fi os.FileInfo, entries []schema.BlobRef) *schema.Schema {
	return &schema.Schema{
		Version:        1,
		Type:           schema.ContentTypeDir,
		FileName:       fi.Name(),
		UnixPermission: fmt.Sprintf("%#o", fi.Mode()),
		DirEntries:     entries,
	}
}

func getFileSchema(fi os.FileInfo, contentRef schema.BlobRef) *schema.Schema {
	return &schema.Schema{
		Version:        1,
		Type:           schema.ContentTypeFile,
		FileName:       fi.Name(),
		UnixPermission: fmt.Sprintf("%#o", fi.Mode()),
		FileParts: []*schema.BytesPart{&schema.BytesPart{
			Size:       uint64(fi.Size()),
			Offset:     0,
			ContentRef: contentRef,
		}},
	}
}
