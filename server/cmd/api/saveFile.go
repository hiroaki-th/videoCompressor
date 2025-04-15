package api

import (
	"encoding/binary"
	"encoding/json"
	"os"
)

const path string = "./tmp/"

type saveFileJson struct {
	Name string `json:"name"`
}

func SaveFile(buff []byte) error {
	jsonSize := int(binary.BigEndian.Uint16(buff[:2]))
	mediaTypeSize := int(int8(buff[2]))
	payloadSize := int(binary.BigEndian.Uint64(buff[3:11]))

	jsonBin := buff[11 : 11+jsonSize]
	_ = buff[11+jsonSize : 11+jsonSize+mediaTypeSize]
	payload := buff[11+jsonSize+mediaTypeSize : 11+jsonSize+mediaTypeSize+payloadSize]

	jsonData := saveFileJson{}
	err := json.Unmarshal(jsonBin, &jsonData)
	if err != nil {
		return err
	}

	file, err := os.Create(path + jsonData.Name)
	if err != nil {
		return err
	}

	_, err = file.Write(payload)
	if err != nil {
		return err
	}

	return nil
}
