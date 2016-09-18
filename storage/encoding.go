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

func (e *encoder) writeTags(arr []tagIdType) {
	e.write(uint8(len(arr)))
	for _, s := range arr {
		e.write(s)
	}
}

func (e *encoder) writeFile(f *StorageFile) {
	e.write(f.Id)
	e.writeString(f.Name)
	e.writeString(f.Period)
	e.writeTags(f.Tags)
	e.write(f.CreationTimestamp)
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

func (d *decoder) readUint8() uint8 {
	var x uint8
	err := binary.Read(d.r, byteOrder, &x)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	return x
}

func (d *decoder) readString() string {
	size := d.readUint8()
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

func (d *decoder) readTags() []tagIdType {
	size := d.readUint8()
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
	s.Tags = d.readTags()
	s.CreationTimestamp = d.readUint64()
	return &s
}
