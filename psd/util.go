package psd

import "os"

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
