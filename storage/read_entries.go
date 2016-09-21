package storage

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

type fileFilter func(*StorageFile, *storageContent) bool

func fileFilterAll(f *StorageFile, content *storageContent) bool {
	return true
}

func (s *Storage) readEntries(filter fileFilter) *storageContent {
	f, err := os.Open(s.metaFilePath())
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}
	defer f.Close()

	return s.readEntriesWithReader(f, filter)
}

func (s *Storage) openAndReadEntries() (*storageContent, *os.File) {
	f, err := os.OpenFile(s.metaFilePath(), os.O_RDWR, 0600)
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}

	return s.readEntriesWithReader(f, nil), f
}

func (s *Storage) readEntriesWithReader(r io.Reader, filter fileFilter) *storageContent {
	content := newContent()

	dec := &decoder{r}
	if dec.readUint64() != magick {
		log.Fatalln("Bad magick")
	}

	for {
		var op uint8
		err := binary.Read(r, byteOrder, &op)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("ERROR", err)
		}

		switch op {
		case opSetAllowedTags:
			content.tags = dec.readTagsData()

		case opSetFile:
			file := dec.readFile()

			if file.Id > content.files.maxFileId {
				content.files.maxFileId = file.Id
			}

			if filter != nil && filter(file, content) {
				content.files.files[file.Id] = file
			}

		case opRemoveFile:
			id := dec.readUint64()
			delete(content.files.files, id)

		default:
			log.Fatalln("Unkwnown op:", op)
		}
	}

	return content
}
