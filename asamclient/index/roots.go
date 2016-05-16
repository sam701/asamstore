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
	roots    map[string]schema.BlobRef
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

func (i *Index) GetRoots() map[string]schema.BlobRef {
	if i.roots == nil {
		i.roots = map[string]schema.BlobRef{}
		f, err := os.Open(i.rootFilePath())
		if err != nil {
			return i.roots
		}
		defer f.Close()

		s := bufio.NewScanner(f)
		for s.Scan() {
			pp := strings.SplitN(s.Text(), " ", 2)
			i.roots[pp[1]] = schema.BlobRef(pp[0])
		}
	}
	return i.roots
}

func (i *Index) GetRootRef(name string) schema.BlobRef {
	return i.GetRoots()[name]
}
