package files

import (
	"errors"

	"github.com/sam701/file-tagger/storage"
	"github.com/urfave/cli"
)

func Add(c *cli.Context) error {
	tags := c.StringSlice("tag")
	period := c.String("period")

	if period == "" {
		return errors.New("Period may not be empty")
	}

	st := storage.Open(c)
	return st.AddFiles(c.Args(), period, tags)
}
