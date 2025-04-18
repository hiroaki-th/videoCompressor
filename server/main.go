package main

import (
	"log"
	"os"
	"videoCompressorServer/cmd/api"
)

func main() {
	server := api.NewServer("tcp", ":8080")
	err := server.Start()
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
