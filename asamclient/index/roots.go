package index

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/sam701/asamstore/asamclient/schema"
)

type Index struct {
	indexDir string
}

func OpenIndex(indexDir string) *Index {
	err := os.MkdirAll(indexDir, 0700)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return &Index{
		indexDir: indexDir,
	}
}

func (i *Index) rootFilePath() string {
	return path.Join(i.indexDir, "roots.txt")
}

func (i *Index) AddRoot(name string, ref schema.BlobRef) {
	ex := i.GetRootRef(name)
	if ex != "" {
		return
	}
	f, err := os.OpenFile(i.rootFilePath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "%s %s\n", ref, name)
}

func (i *Index) GetRootRef(name string) schema.BlobRef {
	f, err := os.Open(i.rootFilePath())
	if err != nil {
		return ""
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		pp := strings.SplitN(s.Text(), " ", 2)
		if pp[1] == name {
			return schema.BlobRef(pp[0])
		}
	}
	return ""
}
