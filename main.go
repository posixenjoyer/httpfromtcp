package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	_ "sync"
)

func processFile(file io.ReadCloser, processChunk func(buffer [8]byte, prefix string, n int) string) <-chan string {
	var fullLine string
	var fileBytes [8]byte

	lineChan := make(chan string)

	go func() {
		for {
			count, err := file.Read(fileBytes[:])
			if err != nil {
				if err != io.EOF {
					fmt.Printf("File Read Error: %v\n", err)
					os.Exit(1)
				}
				if len(fullLine) > 0 {
					lineChan <- fullLine
				}
				break
			}
			fullLine = processChunk(fileBytes, fullLine, count)
		}
		close(lineChan)
	}()

	return lineChan
}

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Println("Failed to open ./messages.txt")
		os.Exit(1)
	}

	processChunk := func(buffer [8]byte, prefix string, count int) string {
		var linePart string

		workingLine := string(buffer[:count])
		parts := strings.Split(workingLine, "\n")

		if len(prefix) > 0 {
			linePart = prefix
		}

		for i, line := range parts {
			if i < len(parts)-1 {
				if len(linePart) > 0 {
					line = linePart + line
					linePart = ""
				}
				fmt.Printf("read: %s\n", line)
				continue
			}

			linePart += line
		}
		return linePart
	}

	fileChan := processFile(file, processChunk)

	for line := range fileChan {
		fmt.Printf("read: %s\n", line)
	}

	if err != nil {
		fmt.Println("We should have never gotten here.")
		os.Exit(1)
	}
	os.Exit(0)
}
