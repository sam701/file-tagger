package storage

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"path"

	"github.com/urfave/cli"
)

type tagIdType uint16

type Storage struct {
	storagePath string
	allowedTags map[string]tagIdType
	tagNames    map[tagIdType]string
	files       map[uint64]*StorageFile

	maxFileId uint64
	maxTagId  tagIdType

	metaFile *os.File
}

func (s *Storage) Close() {
	s.metaFile.Close()
}

const magick uint64 = 0x2132430121324301

const (
	opAddAllowedTag uint8 = iota
	opRemoveAllowedTag
	opAddFile
	opRemoveFile
	opAddTagToFile
	opRemoveTagFromFile
	opRenameTag
)

type StorageFile struct {
	Id                uint64
	Name              string
	Period            string
	Tags              []tagIdType
	CreationTimestamp uint64
}

func (f *StorageFile) removeTag(tag tagIdType) {
	for i, t := range f.Tags {
		if t == tag {
			f.Tags[i] = f.Tags[len(f.Tags)-1]
			f.Tags = f.Tags[:len(f.Tags)-1]
		}
	}
}

func (f *StorageFile) Match(tags []string, tagIdMap map[string]tagIdType) bool {
	if len(f.Tags) < len(tags) {
		return false
	}
	for _, ta := range tags {
		match := false
		tagId := tagIdMap[ta]
		for _, t := range f.Tags {
			if t == tagId {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	return true
}

func (s *Storage) ContentPath(contentPath string) string {
	return path.Join(s.storagePath, contentPath)
}

func (s *Storage) metaFilePath() string {
	return s.ContentPath("meta.bin")
}

func Open(c *cli.Context) *Storage {
	s := &Storage{
		storagePath: c.GlobalString("storage-dir"),
		allowedTags: map[string]tagIdType{},
		tagNames:    map[tagIdType]string{},
		files:       map[uint64]*StorageFile{},
	}
	if _, err := os.Stat(s.storagePath); err != nil {
		log.Fatalln("Storage dir", s.storagePath, "does not exist")
	}

	var err error
	s.metaFile, err = os.OpenFile(s.metaFilePath(), os.O_RDWR, 0600)
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}

	s.checkMagick()
	s.readEntries()
	return s
}

func (s *Storage) readEntries() {
	dec := &decoder{s.metaFile}
	for {
		var op uint8
		err := binary.Read(s.metaFile, byteOrder, &op)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("ERROR", err)
		}

		switch op {
		case opAddAllowedTag:
			tagId := dec.readTagIdType()
			tag := dec.readString()
			s.allowedTags[tag] = tagId
			s.tagNames[tagId] = tag
			s.maxTagId = tagId
		case opRemoveAllowedTag:
			tag := dec.readString()
			delete(s.allowedTags, tag)
		case opAddFile:
			file := dec.readFile()
			s.files[file.Id] = file
			s.maxFileId = file.Id
		case opRemoveFile:
			id := dec.readUint64()
			delete(s.files, id)
		case opAddTagToFile:
			fileId := dec.readUint64()
			tagId := dec.readTagIdType()
			file := s.files[fileId]
			file.Tags = append(file.Tags, tagId)
		case opRemoveTagFromFile:
			fileId := dec.readUint64()
			tagId := dec.readTagIdType()
			s.files[fileId].removeTag(tagId)
		default:
			log.Fatalln("Unkwnown op:", op)
		}
	}
}

func (s *Storage) checkMagick() {
	if (&decoder{s.metaFile}).readUint64() != magick {
		log.Fatalln("Bad magick")
	}
}
