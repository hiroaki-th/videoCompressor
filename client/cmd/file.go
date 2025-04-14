package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func SelectFile(reader *bufio.Reader) (*os.File, error) {
	fmt.Println("please input filename to send server")
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
