package index

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/sam701/asamstore/asamclient/schema"
)

type Commit struct {
	Root       schema.BlobRef            `json:"-"`
	Commit     schema.BlobRef            `json:"commit"`
	Content    schema.BlobRef            `json:"content"`
	CommitTime string                    `json:"commitTime"`
	Changes    []*schema.AttributeChange `json:"changes,omitempty"`
}

func (i *Index) AddCommit(c *Commit) {
	f, err := os.OpenFile(i.commitsFile(c.Root), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(c)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
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
		var c Commit
		err = json.Unmarshal([]byte(scanner.Text()), &c)
		if err != nil {
			log.Fatalln("ERROR", err)
		}

		res = append(res, &c)
	}
	return res
}

func (i *Index) commitsFile(root schema.BlobRef) string {
	return path.Join(i.indexDir, fmt.Sprintf("%s.commits", root))
}
