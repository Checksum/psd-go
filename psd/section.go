package psd

import (
	"encoding/binary"
	"fmt"
	"os"
)

// All sections must implement this interface
// Take in the file handle, read the binary data
// as per the defined structure and return the section
// information
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

func (s Section) String() string {
	return fmt.Sprintf("Section %s starting at offset %d with length %d bytes",
		string(s.Name), s.Start, s.End)
}

// --------------------------------- File Header --------------------------- //
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
	Data FileHeaderBuffer
	Section
}

func (header *FileHeader) Validate() bool {
	// Magic Word
	if string(header.Data.MagicWord[0:]) == "8BPS" && header.Data.Version == 1 &&
		header.Data.Channels < 56 {
		return true
	}
	return false
}

func (header *FileHeader) Read(file *os.File) Section {
	var buffer FileHeaderBuffer
	err := binary.Read(file, binary.BigEndian, &buffer)
	checkError(err)
	header.Data = buffer
	if !header.Validate() {
		panic("Invalid PSD file")
	}
	section := Section{"FileHeader", 0, binary.Size(buffer)}
	return section
}

// --------------------------------- Color Mode ---------------------------- //
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

func (cm *ColorMode) Read(file *os.File) Section {
	start, _ := file.Seek(0, os.SEEK_CUR)
	var length uint32
	err := binary.Read(file, binary.BigEndian, &length)
	checkError(err)
	// If indexed or duotone, we are just going to skip
	// through the data section of color mode
	if length != 0 {
		file.Seek(int64(length), os.SEEK_CUR)
	}
	end, _ := file.Seek(0, os.SEEK_CUR)
	section := Section{"ColorMode", int(start), int(end) - int(start)}
	return section
}

type ResourceBlock struct {
	Signature [4]byte
	Id        uint16
	Name      []byte // Variable
	Size      uint32 // Size of the data field
	Data      []byte
}
