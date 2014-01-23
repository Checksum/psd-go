package psd

import "os"

// Utilities
func CurrentPos(file *os.File) int64 {
	pos, _ := file.Seek(0, os.SEEK_CUR)
	return pos
}

func Pad2(i int) int {
	num := int(i)
	if num%2 != 0 {
		num += 1
	}
	return num
}
