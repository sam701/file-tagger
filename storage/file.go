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
	if len(tags) == 0 {
		log.Fatalln("Empty tag list is not allowed")
	}
	for _, t := range tags {
		if !s.allowedTags[t] {
			log.Fatalln("Tag", t, "is not allowed")
		}
	}

	if period == "" {
		log.Fatalln("Period may not be empty")
	}

	fileName := filepath.Base(filePath)
	contentStoragePath := filepath.Join(s.ContentPath(period), fileName)
	os.MkdirAll(path.Dir(contentStoragePath), 0700)

	copyFile(contentStoragePath, filePath)

	s.maxId++
	s.writeFileEntry(&StorageFile{
		Id:                s.maxId,
		Name:              fileName,
		Period:            period,
		Tags:              tags,
		CreationTimestamp: uint64(time.Now().Unix()),
	})
}

func copyFile(destFileName, srcFileName string) {
	dest, err := os.OpenFile(destFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}
	defer dest.Close()

	src, err := os.Open(srcFileName)
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}
	defer src.Close()

	io.Copy(dest, src)
}

func (s *Storage) writeFileEntry(f *StorageFile) {
	enc := &encoder{s.metaFile}
	enc.write(opAddFile)
	enc.writeFile(f)
}

type FileData struct {
	Id           uint64
	Name         string
	Period       string
	Tags         []string
	CreationDate time.Time
}

func (s *Storage) GetFiles(tags []string) []*FileData {
	out := []*FileData{}
	for _, f := range s.files {
		if f.Match(tags) {
			out = append(out, &FileData{
				Id:           f.Id,
				Name:         f.Name,
				Period:       f.Period,
				Tags:         f.Tags,
				CreationDate: time.Unix(int64(f.CreationTimestamp), 0),
			})
		}
	}
	return out
}
