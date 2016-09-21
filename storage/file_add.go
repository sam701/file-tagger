package storage

import (
	"errors"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

func (s *Storage) AddFiles(filePaths []string, period string, tags []string) error {
	if len(tags) == 0 {
		log.Fatalln("Empty tag list is not allowed")
	}

	content, file := s.openAndReadEntries()
	defer file.Close()

	for _, t := range tags {
		if _, exists := content.tags.allowedTags[t]; !exists {
			return errors.New("Tag " + t + " is not allowed")
		}
	}

	if period == "" {
		return errors.New("Period may not be empty")
	}

	for _, filePath := range filePaths {
		fileName := filepath.Base(filePath)
		contentStoragePath := filepath.Join(s.ContentPath(period), fileName)
		os.MkdirAll(path.Dir(contentStoragePath), 0700)

		copyFile(contentStoragePath, filePath)

		tagIds := make([]tagIdType, len(tags))
		for i, t := range tags {
			tagIds[i] = content.tags.allowedTags[t]
		}

		content.files.maxFileId++

		enc := &encoder{file}
		enc.write(opSetFile)
		enc.writeFile(&StorageFile{
			Id:                content.files.maxFileId,
			Name:              fileName,
			Period:            period,
			Tags:              tagIds,
			CreationTimestamp: uint64(time.Now().Unix()),
		})
	}

	return nil
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
