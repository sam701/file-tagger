package storage

import (
	"fmt"
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
			tx.Rollback()
			log.Fatalln("ERROR", err)
		}
		tx.MustExec("insert into file_tags (file_rowid, tag_rowid) values(?,?)", fileId, tagId)
	}

	tx.Commit()
}

func toIn(vals []string) string {
	out := ""
	for i, v := range vals {
		if i > 0 {
			out += ","
		}
		out += "'" + v + "'"
	}
	return out
}

func (s *Storage) getTagIds(tags []string) []int {
	result := []int{}

	rows, err := s.DB.Query(fmt.Sprintf("select rowid from tags where name in (%s)", toIn(tags)))
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	for rows.Next() {
		var id int
		rows.Scan(&id)
		result = append(result, id)
	}

	return result
}

func (s *Storage) List(tags []string) {
	tagIds := s.getTagIds(tags)

	q := "select f.rowid, f.name, f.period, f.creation_timestamp from files f "
	for i, id := range tagIds {
		if i > 0 {
			q += " and "
		} else {
			q += " where "
		}
		q += fmt.Sprintf(` exists (select 1 from file_tags t 
			where t.file_rowid = f.rowid and t.tag_rowid = %d) `, id)
	}

	rows, err := s.DB.Query(q)

	if err != nil {
		log.Fatalln("ERROR", err)
	}

	for rows.Next() {
		var rowId, ts int
		var name, period string
		err = rows.Scan(&rowId, &name, &period, &ts)
		if err != nil {
			log.Fatalln("ERROR", err)
		}

		fmt.Println(period, name)
	}
}
