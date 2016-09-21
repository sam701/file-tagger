package storage

import (
	"strings"
	"time"

	"github.com/sam701/file-tagger/format"
)

func newFileDate(f *StorageFile, content *storageContent) *format.FileData {
	tagNames := []string{}

	for _, tagId := range f.Tags {
		tagNames = append(tagNames, content.tags.tagNames[tagId])
	}

	return &format.FileData{
		Id:           f.Id,
		Name:         f.Name,
		Period:       f.Period,
		Tags:         tagNames,
		CreationDate: time.Unix(int64(f.CreationTimestamp), 0),
	}
}

func (s *Storage) GetFiles(period string, tags []string) format.Files {
	c := s.readEntries(func(f *StorageFile, content *storageContent) bool {
		return (period == "" || strings.HasPrefix(f.Period, period)) && f.Match(tags, content.tags.allowedTags)
	})

	return c.getFileData()
}
