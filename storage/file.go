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
		if _, exists := s.allowedTags[t]; !exists {
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

	tagIds := make([]tagIdType, len(tags))
	for i, t := range tags {
		tagIds[i] = s.allowedTags[t]
	}

	s.maxFileId++
	s.writeFileEntry(&StorageFile{
		Id:                s.maxFileId,
		Name:              fileName,
		Period:            period,
		Tags:              tagIds,
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
		if f.Match(tags, s.allowedTags) {
			tagNames := make([]string, len(f.Tags))
			for i, tagId := range f.Tags {
				tagNames[i] = s.tagNames[tagId]
			}

			out = append(out, &FileData{
				Id:           f.Id,
				Name:         f.Name,
				Period:       f.Period,
				Tags:         tagNames,
				CreationDate: time.Unix(int64(f.CreationTimestamp), 0),
			})
		}
	}
	return out
}
