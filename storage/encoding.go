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

func (e *encoder) writeStringArray(arr []string) {
	e.write(uint8(len(arr)))
	for _, s := range arr {
		e.writeString(s)
	}
}

func (e *encoder) writeFile(f *StorageFile) {
	e.write(f.Id)
	e.writeString(f.Name)
	e.writeString(f.Period)
	e.writeStringArray(f.Tags)
	e.write(f.CreationTimestamp)
}

type decoder struct {
	r io.Reader
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

func (d *decoder) readStringArray() []string {
	size := d.readUint8()
	out := make([]string, size)

	for ix := 0; ix < int(size); ix++ {
		out[ix] = d.readString()
	}

	return out
}

func (d *decoder) readFile() *StorageFile {
	var s StorageFile
	s.Id = d.readUint64()
	s.Name = d.readString()
	s.Period = d.readString()
	s.Tags = d.readStringArray()
	s.CreationTimestamp = d.readUint64()
	return &s
}
