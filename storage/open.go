package storage

import (
	"log"
	"os"
	"path"

	"github.com/sam701/file-tagger/format"
	"github.com/urfave/cli"
)

type tagIdType uint16

type Storage struct {
	storagePath string
}

type tagsData struct {
	allowedTags map[string]tagIdType
	tagNames    map[tagIdType]string

	maxTagId tagIdType
}

type filesData struct {
	files     map[uint64]*StorageFile
	maxFileId uint64
}

type storageContent struct {
	tags  *tagsData
	files *filesData
}

func newContent() *storageContent {
	return &storageContent{
		tags: &tagsData{
			allowedTags: map[string]tagIdType{},
		},
		files: &filesData{
			files: map[uint64]*StorageFile{},
		},
	}
}

func (sc *storageContent) getFileData() []*format.FileData {
	out := []*format.FileData{}
	for _, f := range sc.files.files {
		out = append(out, newFileDate(f, sc))
	}
	return out
}

const magick uint64 = 0x2132430121324301

const (
	opSetAllowedTags byte = iota
	opSetFile
	opRemoveFile
)

func (s *Storage) ContentPath(contentPath string) string {
	return path.Join(s.storagePath, contentPath)
}

func (s *Storage) metaFilePath() string {
	return s.ContentPath("meta.bin")
}

func Open(c *cli.Context) *Storage {
	s := &Storage{
		storagePath: c.GlobalString("storage-dir"),
	}

	if _, err := os.Stat(s.storagePath); err != nil {
		log.Fatalln("Storage dir", s.storagePath, "does not exist")
	}

	return s
}
