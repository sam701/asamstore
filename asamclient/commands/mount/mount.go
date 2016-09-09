package mount

import (
	"log"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/urfave/cli"
	"github.com/sam701/asamstore/asamclient/client"
	"github.com/sam701/asamstore/asamclient/config"
	"github.com/sam701/asamstore/asamclient/index"
)

var bsClient *client.BlobStorageClient
var bsIndex *index.Index

var userId, groupId uint32
var mountPoint string

func Mount(c *cli.Context) error {
	mountPoint = c.Args().First()
	if mountPoint == "" {
		cli.ShowCommandHelp(c, "mount")
		return nil
	}

	cfg := config.ReadConfig(c.GlobalString("config"))
	bsClient = client.NewClient(cfg)
	bsIndex = index.OpenIndex(cfg.IndexDir)

	readUserAndGroupId()

	var err error
	var conn *fuse.Conn
	conn, err = fuse.Mount(mountPoint)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	defer conn.Close()

	go waitForUnmount()
	err = fs.Serve(conn, &FS{})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	<-conn.Ready
	if err = conn.MountError; err != nil {
		log.Fatalln(err)
	}

	return nil
}

func logUnmountAndExit(args ...interface{}) {
	log.Println(args...)
	fuse.Unmount(mountPoint)
	os.Exit(1)
}

func waitForUnmount() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-c
	log.Println("Unmounting...")
	fuse.Unmount(mountPoint)
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
