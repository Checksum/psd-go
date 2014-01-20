package main

import (
	"./psd"
	"fmt"
)

func main() {
	psd, err := psd.NewPSD("examples/image.psd")
	if err != nil {
		fmt.Println("Error opening file")
		fmt.Println(err)
		return
	}
	psd.Parse()
	psd.File.Close()
}
