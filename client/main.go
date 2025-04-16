package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"videoCompressorClient/cmd"
	"videoCompressorClient/file"
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

	for {

		go func() {
			for {
				message := <-byteMessage
				fmt.Println(len(message))
				_, err := conn.Write(message)
				if err != nil {
					fmt.Println(err)
				}
			}
		}()

		go func() {
			for {
				buff := make([]byte, 0, 1440)
				size, err := conn.Read(buff)
				if err != nil {
					fmt.Println(err)
					y, _ := question("do you want format other file?", reader)
					if y {
						ok = true
					}

				}

				if size > 0 {
					fmt.Println(string(buff))
					y, _ := question("do you want format other file?", reader)
					if y {
						ok = true
					}
				}
			}
		}()

		if ok {
			file, err := file.SelectFile(reader)
			if err != nil {
				fmt.Printf("ERROR: %s", err)
				fmt.Printf("please try again \n\n")
			}
			byteMessage <- cmd.CreateRequest(file)
			ok = false

		}
	}
}
