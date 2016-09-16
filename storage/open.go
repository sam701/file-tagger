package storage

import (
	"log"
	"os"
	"path"

	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli"
)

type Storage struct {
	DB          *sqlx.DB
	storagePath string
}

func (s *Storage) ContentPath(contentPath string) string {
	return path.Join(s.storagePath, contentPath)
}

func (s *Storage) indexDbPath() string {
	return s.ContentPath("index.db")
}

func Open(c *cli.Context) *Storage {
	s := &Storage{
		storagePath: c.GlobalString("storage-dir"),
	}
	if _, err := os.Stat(s.storagePath); err != nil {
		log.Fatalln("Storage dir", s.storagePath, "does not exist")
	}
	db, err := sqlx.Connect("sqlite3", s.indexDbPath())
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	s.DB = db
	return s
}
