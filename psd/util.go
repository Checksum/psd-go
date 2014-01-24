package psd

import (
	"encoding/binary"
	"os"
)

// Pass in a file handle and a reference to a field
// Ex: R(file, &int32)
func R(file *os.File, buf interface{}) error {
	err := binary.Read(file, binary.BigEndian, buf)
	if err != nil {
		panic(err)
	}
	return nil
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func CurrentPos(file *os.File) int64 {
	pos, _ := file.Seek(0, os.SEEK_CUR)
	return pos
}

func Pad2(val interface{}) int {
	i, ok := val.(int)
	if ok {
		i = (i + 1) & ^0x01
	}
	return i
}
