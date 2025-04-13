package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("please input filename to send server")
		filename, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			fmt.Println("please try again")
			continue
		}

		fmt.Println(file)

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
	}
}
