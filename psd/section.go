package psd

import (
	"encoding/binary"
	"fmt"
	"os"
)

type sectionReader interface {
	Read(*os.File) Section
}

// This is common to all sections
// Unfortunately, since binary.Read expects fixed
// length fields, we cannot embed this into every section
// as that will cause the read to try and read the extra bytes
// So for every section, declare the buffer and the section
// separately
type Section struct {
	Name  string
	Start int
	End   int
}

// File Header
type FileHeaderBuffer struct {
	MagicWord [4]byte
	Version   uint16 // 2 bytes
	Reserved  [6]byte
	Channels  uint16
	Height    uint32 // 4 bytes
	Width     uint32
	Depth     uint16
	ColorMode uint16
}

type FileHeader struct {
	Buffer FileHeaderBuffer
	Section
}

func (header *FileHeader) Validate() bool {
	fmt.Println(header)
	// Magic Word
	if string(header.Buffer.MagicWord[0:]) == "8BPS" && header.Buffer.Version == 1 && header.Buffer.Channels < 56 {
		return true
	}
	return false
}

func (header *FileHeader) Read(file *os.File) Section {
	var buffer FileHeaderBuffer
	err := binary.Read(file, binary.BigEndian, &buffer)
	checkError(err)
	header.Buffer = buffer
	if !header.Validate() {
		panic("Invalid PSD file")
	}
	section := Section{"FileHeader", 0, binary.Size(buffer)}
	return section
}

// func (s *FileHeader) Read() {
// 	var header FileHeader
// 	error := binary.Read(s.File, binary.BigEndian, &header)
// 	checkError(error)
// 	if !header.Validate() {
// 		panic("Invalid PSD file")
// 	}
// }

// Color Mode
const (
	CmBitmap       = 0
	CmGrayscale    = 1
	CmIndexed      = 2
	CmRGB          = 3
	CmCMYK         = 4
	CmMultichannel = 7
	CmDuotone      = 8
	CmLab          = 9
)

type ColorMode struct {
	Length uint32
	Data   []byte
}

func (colorMode *ColorMode) Read() {

}

type ResourceBlock struct {
	Signature [4]byte
	Id        uint16
	Name      []byte // Variable
	Size      uint32 // Size of the data field
	Data      []byte
}
