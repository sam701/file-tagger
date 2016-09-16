package main

import (
	"errors"
	"os"

	"github.com/sam701/file-tagger/commands/files"
	"github.com/sam701/file-tagger/commands/tags"
	"github.com/sam701/file-tagger/storage"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "storage-dir, s",
			Usage:  "`PATH` to storage directory",
			EnvVar: "FILE_TAGGER_DIR",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Initialize storage directory",
			Action: storage.Init,
		},
		{
			Name:  "tags",
			Usage: "Manage tags",
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Usage:  "List available tags",
					Action: tags.Print,
				},
				{
					Name:   "add",
					Usage:  "Add allowed tags",
					Action: tags.Add,
				},
			},
		},
		{
			Name:  "files",
			Usage: "Manage files",
			Subcommands: []cli.Command{
				{
					Name:      "add",
					Usage:     "Add files",
					ArgsUsage: "file [file ...]",

					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "period, p",
							Usage: "`PERIOD` to add to the provided files in form year/month",
						},
						cli.StringSliceFlag{
							Name:  "tag, t",
							Usage: "`TAG` to add to the files",
						},
					},
					Action: files.Add,
				},
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		if c.GlobalString("storage-dir") == "" {
			return errors.New("No storage-dir is provided")
		}
		return nil
	}
	app.Run(os.Args)
}
