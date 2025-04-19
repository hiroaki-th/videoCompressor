package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ResponseJson struct {
	Status    uint8     `json:"status"`
	FileName  string    `json:"filename"`
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"timestamp"`
}

func ProcessResponse(buff []byte) error {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	filePath := filepath.Join(homeDir, "Downloads/")
	fmt.Println(filePath)

	jsonSize := int(binary.BigEndian.Uint16(buff[:2]))
	mediaTypeSize := int(int8(buff[2]))
	payloadSize := int(binary.BigEndian.Uint64(buff[3:11]))

	jsonBin := buff[11 : 11+jsonSize]
	jsonData := ResponseJson{}
	err = json.Unmarshal(jsonBin, &jsonData)
	if err != nil {
		return err
	}
	fmt.Printf("{\n status: %d,\n message: %s,\n}\n\n", jsonData.Status, jsonData.Message)

	mediaType := "bin"
	if mediaTypeSize > 0 {
		mediaType = string(buff[11+jsonSize : 11+jsonSize+mediaTypeSize])
	}

	if payloadSize > 0 {
		payload := buff[11+jsonSize+mediaTypeSize : 11+jsonSize+mediaTypeSize+payloadSize]
		file, err := os.Create(filePath + "formattedFile." + mediaType)
		if err != nil {
			return err
		}

		_, err = file.Write(payload)
		if err != nil {
			return err
		}
	}

	fmt.Println("successfully downloaded file. please check /Users[Home]/Downloads/")
	os.Exit(0)
	return nil
}

func GetTotalSize(buff []byte) int {
	jsonSize := int(binary.BigEndian.Uint16(buff[:2]))
	mediaTypeSize := int(int8(buff[2]))
	payloadSize := int(binary.BigEndian.Uint64(buff[3:11]))

	return 11 + jsonSize + mediaTypeSize + payloadSize
}
