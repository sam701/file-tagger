package storage

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

func (s *Storage) AddFile(filePath, period string, tags []string) {
	fileName := filepath.Base(filePath)
	contentStoragePath := filepath.Join(s.ContentPath(period), fileName)
	os.MkdirAll(path.Dir(contentStoragePath), 0700)

	dest, err := os.OpenFile(contentStoragePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}
	defer dest.Close()

	src, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}
	defer src.Close()

	io.Copy(dest, src)

	tx := s.DB.MustBegin()

	relPath, _ := filepath.Rel(s.storagePath, contentStoragePath)
	fileId, _ := tx.MustExec("insert into files (path, name, period, creation_timestamp) values(?,?,?,?)",
		relPath, fileName, period, time.Now().Unix()).LastInsertId()

	for _, tag := range tags {
		var tagId int
		err := tx.Get(&tagId, "select rowid from tags where name = ?", tag)
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		tx.MustExec("insert into file_tags (file_rowid, tag_rowid) values(?,?)", fileId, tagId)
	}

	tx.Commit()
}
