package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)

		currentLine := ""
		buffer := make([]byte, 8)

		for {
			byteRead, err := f.Read(buffer)
			if byteRead > 0 {
				str := string(buffer[:byteRead])
				parts := strings.Split(str, "\n")

				for i := 0; i < len(parts)-1; i++ {
					lines <- currentLine + parts[i]
					currentLine = ""
				}

				currentLine += parts[len(parts)-1]
			}

			if errors.Is(err, io.EOF) {
				break
			}
		}

		if currentLine != "" {
			lines <- currentLine
		}
	}()

	return lines
}

const port = ":42069"

func main() {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())

		lines := getLinesChannel(conn)

		for line := range lines{
			fmt.Println(line)
		}
		
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}