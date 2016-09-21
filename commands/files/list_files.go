package files

import (
	"github.com/sam701/file-tagger/storage"
	"github.com/urfave/cli"
)

func List(c *cli.Context) error {
	files := storage.Open(c).GetFiles(c.String("period"), c.Args())
	files.Print()
	return nil
}
