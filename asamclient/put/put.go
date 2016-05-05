package put

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/sam701/asamstore/asamclient/client"
	"github.com/sam701/asamstore/asamclient/config"
)

func PutAction(c *cli.Context) error {
	cfg := config.ReadConfig(c.GlobalString("config"))
	cl := client.NewClient(cfg)

	cl.Put("1234567891234", strings.NewReader("hello, this is a test"))

	return nil
}
