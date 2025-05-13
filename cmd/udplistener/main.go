package main

import (
	"bufio"
	"fmt"
	_ "io"
	"net"
	"os"
	_ "strings"
	_ "sync"
)

func checkAndPrintError(message string, err error) {
	if err != nil {
		fmt.Printf("%s: %v", message, err)
	}
}

func connectionHandler(udpListener *net.UDPConn) error {
	stdinReader := bufio.NewReader(os.Stdin)

	/*
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
					fmt.Printf("%s\n", line)
					continue
				}

				linePart += line
			}
			return linePart
		}
	*/

	for {
		fmt.Print("> ")
		line, err := stdinReader.ReadString('\n')
		checkAndPrintError("Error reading from stdin: ", err)

		_, err = udpListener.Write([]byte(line))
		checkAndPrintError("UDP Write failed: ", err)

		/*
			httpChan := processConnections(httpConn, processChunk)

			for line := range httpChan {
				fmt.Printf("%s\n", line)
			}

			fmt.Printf("Connection from %s closed.\n", httpPeer)
		*/
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
	udpRemoteAddr, err := net.ResolveUDPAddr("udp", remoteAddress)
	if err != nil {
		fmt.Println("Failed to open ./messages.txt")
		os.Exit(1)
	}

	udpSession, err := net.DialUDP("udp", nil, udpRemoteAddr)
	checkAndPrintError("UDP Dial Failed", err)

	defer udpSession.Close()

	err = connectionHandler(udpSession)
	if err != nil {
		fmt.Printf("Error handling connection: %v\n", err)
	}
}
