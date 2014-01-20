package psd

import (
	"encoding/binary"
	"log"
	"os"
)

type PSD struct {
	Path   string
	File   *os.File `json:"-"` // Exclude from JSON
	Width  uint32   `json:"width"`
	Height uint32   `json:"height"`
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func NewPSD(path string) (*PSD, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	psd := &PSD{
		Path: path,
		File: file,
	}
	return psd, nil
}

func (psd *PSD) Parse() {

	sections := []sectionReader{new(FileHeader)}
	for _, section := range sections {
		section.Read(psd.File)
	}

	// psd.ParseHeader()
	// // Skip through the color mode
	// psd.ParseColorMode()
	// // Parse the resource blocks
	// psd.ParseResourceBlocks()
	// // Parse Layer info
	// psd.ParseLayerInfo()
	// //
	// psd.ParseImageData()
}

func (psd *PSD) ParseHeader() {
	var header FileHeader
	error := binary.Read(psd.File, binary.BigEndian, &header)
	checkError(error)
	if !header.Validate() {
		panic("Invalid PSD file")
	}
}

// For now, we are not going to bother with variable
// color modes (indexed and variable)
func (psd *PSD) ParseColorMode() {
	var cm uint32
	error := binary.Read(psd.File, binary.BigEndian, &cm)
	checkError(error)
	log.Printf("ColorMode: %d", cm)
	// If indexed or duotone, we are just going to skip
	// through the data section of color mode
	if cm == CmIndexed || cm == CmDuotone {
		psd.File.Seek(int64(cm), os.SEEK_CUR)
	}
}

// Parse Resource Blocks
func (psd *PSD) ParseResourceBlocks() {
	var Len uint32
	error := binary.Read(psd.File, binary.BigEndian, &Len)
	checkError(error)
	log.Printf("Length of Image Resource Section: %d bytes", Len)

	psd.File.Seek(int64(Len), os.SEEK_CUR)

	// // Signature
	// var Signature [4]byte
	// binary.Read(psd.File, binary.BigEndian, &Signature)
	// log.Printf(string(Signature[0:]))
	// // Resource ID
	// var Id uint16
	// binary.Read(psd.File, binary.BigEndian, &Id)
	// log.Printf("Resource ID: %d", Id)
	// // Resource Name
	// var Name [2]byte
	// binary.Read(psd.File, binary.BigEndian, &Name)
	// fmt.Println(Name)
	// log.Printf("Resource Name: %s", string(Name[:]))
	// // Length of Data
	// var Size uint32
	// binary.Read(psd.File, binary.BigEndian, &Size)
	// if Size % 2 != 0 {
	// 	Size += 1
	// }
	// log.Printf("Size of data segment: %d", Size)

	// psd.File.Seek(int64(Size), os.SEEK_CUR)
}

// Parse LayerInfo
func (psd *PSD) ParseLayerInfo() {
	var Len uint32
	error := binary.Read(psd.File, binary.BigEndian, &Len)
	checkError(error)
	log.Printf("Length of Layer info section %d", Len)
	psd.File.Seek(int64(Len), os.SEEK_CUR)
}

func (psd *PSD) ParseImageData() {
	var Len uint16
	binary.Read(psd.File, binary.BigEndian, &Len)
	log.Printf("Length of Image Data: %d", Len)
}

func (psd *PSD) CurrentPos() int64 {
	pos, _ := psd.File.Seek(0, os.SEEK_CUR)
	return pos
}
