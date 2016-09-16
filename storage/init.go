package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli"
)

func Init(c *cli.Context) error {
	storage := &Storage{
		storagePath: c.GlobalString("storage-dir"),
	}

	if _, err := os.Stat(storage.indexDbPath()); err == nil {
		fmt.Println("Storage", storage.storagePath, "was already initialized")
		return nil
	}

	err := os.Mkdir(storage.storagePath, 0700)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	db, err := sqlx.Connect("sqlite3", storage.indexDbPath())
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	defer db.Close()

	db.MustExec(`
    create table tags(
        name varchar(128) primary key
    );

    create table files(
        path varchar(256) primary key,
		name varchar(64),
		period varchar(32),
        creation_timestamp integer
    );

    create table file_tags(
        file_rowid integer,
        tag_rowid integer
    );

    create index tags_ix on file_tags(tag_rowid);
    `)

	fmt.Println("Storage has been initialsed:", storage.storagePath)

	return nil
}
