package api

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
)

const path string = "../../tmp/"

type saveFileJson struct {
	FileName string
}

func SaveFile(buff []byte) error {

	jsonSize := binary.BigEndian.Uint64(buff[:5])
	mediaTypeSize := binary.BigEndian.Uint64(buff[5:6])
	payloadSize := binary.BigEndian.Uint64(buff[6:7])

	jsonBin := buff[7:jsonSize]
	mediaType := buff[7+jsonSize : 7+jsonSize+mediaTypeSize]
	payload := buff[7+jsonSize+mediaTypeSize : 7+jsonSize+mediaTypeSize+payloadSize]

	jsonData := saveFileJson{}
	err := json.Unmarshal(jsonBin, &jsonData)
	if err != nil {
		return err
	}

	mediaTypeStr := string(mediaType)
	fmt.Println(mediaTypeStr)

	err = os.WriteFile(path+jsonData.FileName, payload, os.FileMode(os.O_APPEND)|os.FileMode(os.O_CREATE)|os.FileMode(os.O_RDWR))
	if err != nil {
		return err
	}

	return nil
}
