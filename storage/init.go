package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func Init(c *cli.Context) error {
	storage := &Storage{
		storagePath: c.GlobalString("storage-dir"),
	}

	if _, err := os.Stat(storage.metaFilePath()); err == nil {
		fmt.Println("Storage", storage.storagePath, "was already initialized")
		return nil
	}

	err := os.Mkdir(storage.storagePath, 0700)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	f, err := os.OpenFile(storage.metaFilePath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}
	defer f.Close()

	(&encoder{f}).write(magick)

	fmt.Println("Storage has been initialsed:", storage.storagePath)

	return nil
}
