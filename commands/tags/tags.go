package tags

import (
	"fmt"

	"github.com/sam701/file-tagger/storage"
	"github.com/urfave/cli"
)

func Print(c *cli.Context) error {
	st := storage.Open(c)
	rows, _ := st.DB.Query("select name from tags order by name")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Println(name)
	}
	return nil
}

func Add(c *cli.Context) error {
	st := storage.Open(c)
	defer st.DB.Close()

	for _, tag := range c.Args() {
		st.DB.MustExec("insert into tags (name) values(?)", tag)
		fmt.Println("Inserted tag:", tag)
	}

	return nil
}
