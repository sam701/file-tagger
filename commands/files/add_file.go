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
	defer st.Close()

	for _, file := range c.Args() {
		st.AddFile(file, period, tags)
	}

	return nil
}
