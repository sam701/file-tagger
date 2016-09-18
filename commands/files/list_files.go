package files

import (
	"fmt"

	"github.com/sam701/file-tagger/storage"
	"github.com/urfave/cli"
)

func List(c *cli.Context) error {
	for _, f := range storage.Open(c).GetFiles(c.Args()) {
		fmt.Println(f.Id, f.Name, f.Period, f.Tags)
	}
	return nil
}
