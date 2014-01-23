package psd

import (
	"encoding/binary"
	"os"
)

// Pascal length prefixed string
// The first byte is the length of the string
// This is even padded
type pString struct {
	Len  uint8
	Text []byte
}

func (str *pString) Read(file *os.File) {
	var length uint8
	err := binary.Read(file, binary.BigEndian, &length)
	checkError(err)
	// fmt.Println(length)
	str.Len = length
	// Pad the length
	// Note that the even padded string size also includes
	// the first byte that is read as length
	length += 1
	if length%2 != 0 {
		length += 1
	}
	// Since we have already read one byte from the even padded
	// string, decrease the number of bytes to read by 1
	length -= 1
	// And now, read length bytes as the string
	str.Text = make([]byte, length)
	// fmt.Println("Size of text: ", len(str.Text))
	e := binary.Read(file, binary.BigEndian, &str.Text)
	checkError(e)
}

// Unicode string
type uString struct {
	Len  uint32 // Number of characters (not bytes)
	Text []byte // 2 bytes per character
}

func (str *uString) Read(file *os.File) {
	var length uint32
	err := binary.Read(file, binary.BigEndian, &length)
	checkError(err)
	str.Len = length
	str.Text = make([]byte, length*2)
	e := binary.Read(file, binary.BigEndian, &str.Text)
	checkError(e)
}
