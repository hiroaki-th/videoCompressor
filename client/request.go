package main

import (
	"os"
	"strings"
)

func CreateRequest(file *os.File) []byte {

	filename := []byte(file.Name())
	mediaType := []byte(strings.Split(file.Name(), ".")[1])

	fileBody := make([]byte, 0)
	_, err := file.Read(fileBody)
	if err != nil {
		return nil
	}

	request := make([]byte, 0)
	return setInfo(request, filename, mediaType, fileBody)
}

func setInfo(base []byte, params ...[]byte) []byte {
	for _, value := range params {
		base = append(base, value...)
	}

	return base
}
