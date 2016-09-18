package tags

import (
	"fmt"

	"github.com/sam701/file-tagger/storage"
	"github.com/urfave/cli"
)

func Print(c *cli.Context) error {
	st := storage.Open(c)
	for _, t := range st.GetTags() {
		fmt.Println(t)
	}
	return nil
}

func Add(c *cli.Context) error {
	st := storage.Open(c)
	defer st.Close()

	for _, tag := range c.Args() {
		st.AddTag(tag)
		fmt.Println("Inserted tag:", tag)
	}

	return nil
}
