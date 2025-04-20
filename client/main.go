package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"videoCompressorClient/cmd"
	myfile "videoCompressorClient/file"
)

func question(message string, reader *bufio.Reader) (bool, error) {
	fmt.Printf("%s [y/n]", message)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("unexpected error has occurred", err)
		os.Exit(-1)
	}

	if strings.Contains(input, "n") {
		fmt.Println("see you again!")
		os.Exit(0)
	}

	return true, nil
}

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	byteMessage := make(chan []byte)

	var ok bool = true

	// send request to server
	go func() {
		for {
			message := <-byteMessage
			_, err := conn.Write(message)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	// process response from server
	go func() {
		response := make([]byte, 0)
		totalSize := 0
		for {
			buff := make([]byte, 1440)
			size, err := conn.Read(buff)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			if size > 0 {
				response = append(response, buff[:size]...)
			}

			if len(response) > 11 {
				if totalSize == 0 {
					totalSize = cmd.GetTotalSize(response)
				}

				if totalSize == len(response) {
					cmd.ProcessResponse(response)
				}
			}

		}
	}()

	// select file
	for {
		if ok {
			file, err := myfile.SelectFile(reader)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
				fmt.Printf("please try again \n\n")
				continue
			}

			json, err := myfile.SelectFormat(reader, file)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
				fmt.Printf("please try again \n\n")
				continue
			}

			byteMessage <- cmd.CreateRequest(file, json)
			ok = false
		}
	}
}
