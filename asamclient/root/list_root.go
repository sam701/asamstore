package root

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/sam701/asamstore/asamclient/config"
	"github.com/sam701/asamstore/asamclient/index"
)

func listRoots(c *cli.Context) {
	conf := config.ReadConfig(c.GlobalString("config"))
	ix := index.OpenIndex(conf.IndexDir)

	for name, hash := range ix.GetRoots() {
		fmt.Println(hash, name)
	}
}
