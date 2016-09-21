package storage

import (
	"encoding/binary"
	"io"
	"log"
)

type encoder struct {
	w io.Writer
}

var byteOrder = binary.BigEndian

func (e *encoder) write(v interface{}) {
	err := binary.Write(e.w, byteOrder, v)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
}

func (e *encoder) writeString(s string) {
	data := []byte(s)
	e.write(uint8(len(data)))
	err := binary.Write(e.w, byteOrder, data)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
}

func (e *encoder) writeFileTags(arr []tagIdType) {
	e.write(uint8(len(arr)))
	for _, s := range arr {
		e.write(s)
	}
}

func (e *encoder) writeFile(f *StorageFile) {
	e.write(f.Id)
	e.writeString(f.Name)
	e.writeString(f.Period)
	e.writeFileTags(f.Tags)
	e.write(f.CreationTimestamp)
}

func (e *encoder) writeAllowedTagsMap(tags map[string]tagIdType) {
	e.write(byte(len(tags)))
	for name, id := range tags {
		e.writeString(name)
		e.write(id)
	}
}

type decoder struct {
	r io.Reader
}

func (d *decoder) readTagIdType() tagIdType {
	var x tagIdType
	err := binary.Read(d.r, byteOrder, &x)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	return x
}

func (d *decoder) readUint64() uint64 {
	var x uint64
	err := binary.Read(d.r, byteOrder, &x)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	return x
}

func (d *decoder) readByte() byte {
	var x byte
	err := binary.Read(d.r, byteOrder, &x)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	return x
}

func (d *decoder) readString() string {
	size := d.readByte()
	buf := make([]byte, size)

	n, err := d.r.Read(buf)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	if n != int(size) {
		log.Fatalln("Expected", size, "but read", n)
	}
	return string(buf)
}

func (d *decoder) readFileTags() []tagIdType {
	size := d.readByte()
	out := make([]tagIdType, size)

	err := binary.Read(d.r, byteOrder, out)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return out
}

func (d *decoder) readFile() *StorageFile {
	var s StorageFile
	s.Id = d.readUint64()
	s.Name = d.readString()
	s.Period = d.readString()
	s.Tags = d.readFileTags()
	s.CreationTimestamp = d.readUint64()
	return &s
}

func (d *decoder) readTagsData() *tagsData {
	td := &tagsData{
		allowedTags: map[string]tagIdType{},
		tagNames:    map[tagIdType]string{},
	}
	size := int(d.readByte())
	for i := 0; i < size; i++ {
		name := d.readString()
		id := d.readTagIdType()

		td.allowedTags[name] = id
		td.tagNames[id] = name
		if id > td.maxTagId {
			td.maxTagId = id
		}
	}
	return td
}
