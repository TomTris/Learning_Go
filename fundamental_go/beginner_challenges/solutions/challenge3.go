package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

const BUFFER_SIZE = 1

func main() {
	buffer := make([]byte, BUFFER_SIZE)
	f, err := os.Open("../README.md")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	total_bytes := 0
	content := ""

	for {
		read_byte, err := f.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		content += string(buffer[:read_byte])
		total_bytes += read_byte
	}
	fmt.Println(content)
	fmt.Print(total_bytes)
}
