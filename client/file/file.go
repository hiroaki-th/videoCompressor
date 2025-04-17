package file

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type FileJson struct {
	Name   string `json:"name"`
	Format string `json:"format"`
}

func SelectFile(reader *bufio.Reader) (*os.File, error) {
	fmt.Printf("please input filename to send server\n\n")
	filename, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	filename = strings.Trim(filename, "\n")

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}
