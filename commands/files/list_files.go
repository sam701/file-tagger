package files

import (
	"github.com/sam701/file-tagger/storage"
	"github.com/urfave/cli"
)

func List(c *cli.Context) error {
	storage.Open(c).List(c.Args())
	return nil
}
