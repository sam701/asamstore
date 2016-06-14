package put

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/codegangsta/cli"
	"github.com/sam701/asamstore/asamclient/client"
	"github.com/sam701/asamstore/asamclient/config"
	"github.com/sam701/asamstore/asamclient/index"
	"github.com/sam701/asamstore/asamclient/schema"
)

var bsClient *client.BlobStorageClient

func PutAction(c *cli.Context) error {
	cfg := config.ReadConfig(c.GlobalString("config"))
	bsClient = client.NewClient(cfg)

	contentPath := c.Args().First()
	if contentPath == "" {
		log.Fatalln("No content path provided")
	}
	var err error
	contentPath, err = filepath.Abs(contentPath)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	rootName := c.String("root")
	if rootName == "" {
		log.Fatalln("No root given")
	}

	ix := index.OpenIndex(cfg.IndexDir)
	rootRef := ix.GetRootRef(rootName)
	if rootRef == "" {
		log.Fatalln("No such root", rootName)
	}

	ref := putFile(contentPath)
	commits := ix.GetCommits(rootRef)
	ll := len(commits)
	if ll > 0 && commits[ll-1].Content == ref {
		log.Println("No changes")
	} else {
		cs := getCommitSchema(rootRef, ref)
		commitRef := bsClient.PutSchema(cs)
		ix.AddCommit(&index.Commit{rootRef, commitRef, ref, cs.CommitTime})
		bsClient.UpdateServerState()
	}

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

		sort.Sort(byName(fis))
		for _, fi := range fis {
			ref := putFile(path.Join(filePath, fi.Name()))
			entries = append(entries, ref)
		}
		return bsClient.PutSchema(getDirSchema(st, entries))
	} else {
		return putFileParts(f, st)
	}
}

type byName []os.FileInfo

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

func getDirSchema(fi os.FileInfo, entries []schema.BlobRef) *schema.Schema {
	s := schema.NewSchema(schema.ContentTypeDir)
	s.FileName = fi.Name()
	s.UnixPermission = fmt.Sprintf("%#o", fi.Mode())
	s.UnixMtime = fi.ModTime().Format(time.RFC3339)
	s.DirEntries = entries
	return s
}

func getCommitSchema(root, content schema.BlobRef) *schema.Schema {
	s := schema.NewSchema(schema.ContentTypeCommit)
	s.RootRef = root
	s.CommitTime = time.Now().Format(time.RFC3339)
	s.ContentRef = content
	return s
}
