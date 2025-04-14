package api

import (
	"encoding/json"
	"fmt"
	"os"
)

const path string = "../../tmp/"

type saveFileJson struct {
	Name string `json:"name"`
}

func SaveFile(buff []byte) error {
	jsonSize := int(buff[0])
	mediaTypeSize := int(buff[1])
	payloadSize := int(buff[2])

	jsonBin := buff[3:jsonSize]
	mediaType := buff[3+jsonSize : 3+jsonSize+mediaTypeSize]
	payload := buff[3+jsonSize+mediaTypeSize : 3+jsonSize+mediaTypeSize+payloadSize]

	jsonData := saveFileJson{}
	err := json.Unmarshal(jsonBin, &jsonData)
	if err != nil {
		return err
	}

	mediaTypeStr := string(mediaType)
	fmt.Println(mediaTypeStr)

	err = os.WriteFile(path+jsonData.Name, payload, os.FileMode(os.O_APPEND)|os.FileMode(os.O_CREATE)|os.FileMode(os.O_RDWR))
	if err != nil {
		return err
	}

	return nil
}
