package mount

import (
	"log"
	"os/user"
	"strconv"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/codegangsta/cli"
	"github.com/sam701/asamstore/asamclient/client"
	"github.com/sam701/asamstore/asamclient/config"
	"github.com/sam701/asamstore/asamclient/index"
)

var bsClient *client.BlobStorageClient
var bsIndex *index.Index

var userId, groupId uint32

func Mount(c *cli.Context) error {
	mountPoint := c.Args().First()
	if mountPoint == "" {
		cli.ShowCommandHelp(c, "mount")
		return nil
	}

	cfg := config.ReadConfig(c.GlobalString("config"))
	bsClient = client.NewClient(cfg)
	bsIndex = index.OpenIndex(cfg.IndexDir)

	readUserAndGroupId()

	ctx, err := fuse.Mount(mountPoint)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	defer ctx.Close()

	err = fs.Serve(ctx, &FS{})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	<-ctx.Ready
	if err := ctx.MountError; err != nil {
		log.Fatalln(err)
	}

	return nil
}

func readUserAndGroupId() {
	u, err := user.Current()
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	ui, err := strconv.Atoi(u.Uid)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	userId = uint32(ui)

	gi, err := strconv.Atoi(u.Gid)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	groupId = uint32(gi)
}

type FS struct{}

func (f *FS) Root() (fs.Node, error) {
	return &storageNode{}, nil
}
