package main

import (
	"fmt"
	"github.com/posixenjoyer/httpfromtcp/internal/requests"
	"net"
	"os"
	_ "strings"
	_ "sync"
)

func connectionHandler(httpListener net.Listener) error {
	for {
		httpConn, err := httpListener.Accept()
		if err != nil {
			fmt.Printf("Http connection Accept() failure: %v\n", err)
			return err
		}

		httpPeer := httpConn.RemoteAddr().String()

		fmt.Printf("Accepted HTTP Connection: %s\n", httpPeer)
		httpRequest, err := request.RequestFromReader(httpConn)

		if !httpRequest.RequestLine.IsEmpty() {
			fmt.Printf("Request line:\n")
			fmt.Printf("- Method: %s\n", httpRequest.RequestLine.Method)
			fmt.Printf("- Target: %s\n", httpRequest.RequestLine.RequestTarget)
			fmt.Printf("- Version: %s\n", httpRequest.RequestLine.HttpVersion)
		}

		fmt.Printf("Connection from %s closed.\n", httpPeer)
	}
}

/*
func processConnections(file io.ReadCloser, processChunk func(buffer [8]byte, prefix string, n int) string) <-chan string {
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
*/

func main() {
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		fmt.Println("Failed to open ./messages.txt")
		os.Exit(1)
	}

	defer listener.Close()

	err = connectionHandler(listener)
	if err != nil {
		fmt.Printf("Error handling connection: %v\n", err)
	}
}
