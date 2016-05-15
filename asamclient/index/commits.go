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

type Commit struct {
	Root       schema.BlobRef
	Commit     schema.BlobRef
	Content    schema.BlobRef
	CommitTime string
}

func (i *Index) AddCommit(c *Commit) {
	f, err := os.OpenFile(i.commitsFile(c.Root), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}
	defer f.Close()

	fmt.Fprintln(f, c.Commit, c.Content, c.CommitTime)
}

func (i *Index) GetCommits(root schema.BlobRef) []*Commit {
	res := []*Commit{}
	f, err := os.Open(i.commitsFile(root))
	if err != nil {
		return res
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		pp := strings.Split(scanner.Text(), " ")
		res = append(res, &Commit{
			Root:       root,
			Commit:     schema.BlobRef(pp[0]),
			Content:    schema.BlobRef(pp[1]),
			CommitTime: pp[2],
		})
	}
	return res
}

func (i *Index) commitsFile(root schema.BlobRef) string {
	return path.Join(i.indexDir, fmt.Sprintf("%s.commits", root))
}
