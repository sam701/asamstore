package ls

import (
	"fmt"
	"log"
	"strings"

	"github.com/sam701/asamstore/asamclient/client"
	"github.com/sam701/asamstore/asamclient/config"
	"github.com/sam701/asamstore/asamclient/index"
	"github.com/sam701/asamstore/asamclient/schema"
	"github.com/urfave/cli"
)

var bsClient *client.BlobStorageClient

func LsAction(c *cli.Context) error {
	cfg := config.ReadConfig(c.GlobalString("config"))
	bsClient = client.NewClient(cfg)
	ix := index.OpenIndex(cfg.IndexDir)

	rootName := c.String("root")
	if rootName == "" {
		log.Fatalln("No root provided")
	}

	requiredTags := c.Args()
	if requiredTags == nil || len(requiredTags) == 0 {
		log.Fatalln("No tags provided")
	}

	rootRef := ix.GetRoots()[rootName]
	commits := ix.GetCommits(rootRef)

	for _, commit := range commits {
		changes := commit.Changes
		if changes != nil {
			for _, change := range changes {
				if change.AttributeName == "tags" && hasAllTags(requiredTags, change.Values) {
					printCommit(commit.Commit, commit.Content, requiredTags)
				}
			}
		}
	}

	return nil
}

func printCommit(commitRef, contentRef schema.BlobRef, requiredTags []string) {
	commit := bsClient.GetSchema(commitRef)
	content := bsClient.GetSchema(contentRef)

	fmt.Println(contentRef, content.FileName, commit.ContentAttributeChanges[0].Values)
}

func hasAllTags(requiredTags []string, availableTags []string) bool {
	for _, rt := range requiredTags {
		found := false
		for _, tag := range availableTags {
			if strings.HasPrefix(tag, rt) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
