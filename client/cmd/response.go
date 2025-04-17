package cmd

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"time"
)

type ResponseJson struct {
	Status    uint8     `json:"status"`
	FileName  string    `json:"filename"`
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"timestamp"`
}

const path string = "./tmp/"

func ProcessResponse(buff []byte) error {
	jsonSize := int(binary.BigEndian.Uint16(buff[:2]))
	mediaTypeSize := int(int8(buff[2]))
	payloadSize := int(binary.BigEndian.Uint64(buff[3:11]))

	jsonBin := buff[11 : 11+jsonSize]
	_ = buff[11+jsonSize : 11+jsonSize+mediaTypeSize]
	payload := buff[11+jsonSize+mediaTypeSize : 11+jsonSize+mediaTypeSize+payloadSize]

	jsonData := ResponseJson{}
	err := json.Unmarshal(jsonBin, &jsonData)
	if err != nil {
		return err
	}

	if mediaTypeSize > 0 {
		_ = buff[11+jsonSize : 11+jsonSize+mediaTypeSize]
	}

	if payloadSize > 0 {
		file, err := os.Create(path + jsonData.FileName)
		if err != nil {
			return err
		}

		_, err = file.Write(payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetTotalSize(buff []byte) int {
	jsonSize := int(binary.BigEndian.Uint16(buff[:2]))
	mediaTypeSize := int(int8(buff[2]))
	payloadSize := int(binary.BigEndian.Uint64(buff[3:11]))

	return 11 + jsonSize + mediaTypeSize + payloadSize
}
