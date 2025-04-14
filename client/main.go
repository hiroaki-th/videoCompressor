package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"videoCompressorClient/cmd"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	reader := bufio.NewReader(os.Stdin)
	byteMessage := make(chan []byte)

	for {

		go func() {
			for {
				message := <-byteMessage
				_, err := conn.Write(message)
				if err != nil {
					fmt.Println(err)
				}
			}
		}()

		go func() {
			for {
				buff := make([]byte, 1440)
				size, err := conn.Read(buff)
				if err != nil {
					fmt.Println(err)
					continue
				}

				if size > 0 {
					fmt.Println(string(buff))
				}
			}
		}()

		file, err := cmd.SelectFile(reader)
		if err != nil {
			fmt.Printf("ERROR: %s", err)
			fmt.Printf("please try again \n\n")
			continue
		}

		byteMessage <- cmd.CreateRequest(file)
	}
}
