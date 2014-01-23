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
	start := CurrentPos(file)
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

// ------------------------------- Resource Blocks ------------------------- //
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

type ImageResources struct {
	Length uint32
	Data   []byte
}

type ResourceBlock struct {
	Signature [4]byte
	Id        uint16
	Name      pString // Variable length string
	Size      uint32  // Size of the data field
	Data      []byte
}

func (r *ResourceBlock) Read(file *os.File) int {
	start := CurrentPos(file)
	// The block signature is always 8BIM
	binary.Read(file, binary.BigEndian, &r.Signature)
	fmt.Println(string(r.Signature[0:]))
	// Next comes the ID
	binary.Read(file, binary.BigEndian, &r.Id)
	fmt.Println("Resource ID: ", r.Id)
	// Read the resource name
	// r.Name = new(pString)
	// fmt.Printf("Before reading %d\n", CurrentPos(file))
	r.Name.Read(file)
	// fmt.Printf("After reading %d\n", CurrentPos(file))
	// Block size
	binary.Read(file, binary.BigEndian, &r.Size)
	if r.Size%2 != 0 {
		r.Size += 1
	}
	fmt.Println("Size of data section: ", r.Size)
	r.Data = make([]byte, r.Size)
	binary.Read(file, binary.BigEndian, &r.Data)
	// fmt.Println(string(r.Data[:]))
	end := CurrentPos(file)
	return int(end - start)
}

// The Resource block is of variable size. To determine
// where we have to stop, we need to first read the length
// of the resource section (which is the first 4 bytes) and then
// loop over and process the individual resource blocks
// Each resource block, in turn can be of variable size
func (b *ImageResources) Read(file *os.File) Section {
	start, _ := file.Seek(0, os.SEEK_CUR)
	var length uint32
	read := uint32(0)
	err := binary.Read(file, binary.BigEndian, &length)
	checkError(err)
	fmt.Printf("\n\nLength of resource block section: %d bytes\n\n", length)

	for read < length {
		fmt.Println("Creating new block at position: ", CurrentPos(file))
		// Create a new resource block
		block := new(ResourceBlock)
		size := block.Read(file)
		fmt.Println("Position after reading block: ", CurrentPos(file))
		read += uint32(size)
		fmt.Printf("Resource Block %s, Id: %d, Size: %d, Pos:%d \n",
			string(block.Name.Text[0:]), block.Id, block.Size, read)
	}

	return Section{"ImageResources", int(start), int(length)}
}
