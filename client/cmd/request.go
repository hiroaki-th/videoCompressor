package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Header struct {
	JsonSize      int16
	MediaTypeSize int8
	PayloadSize   int64
}

type Body struct {
	Json      []byte
	MediaType []byte
	Payload   []byte
}

type FileJson struct {
	Name string `json:"name"`
}

func CreateRequest(file *os.File) []byte {

	header := Header{}
	body := Body{}

	// get json and jsonSize
	filename := string([]byte(file.Name()))
	fileJson := FileJson{
		Name: filename,
	}
	byteJson, err := json.Marshal(fileJson)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	header.JsonSize = int16(len(byteJson))
	body.Json = byteJson

	// get mediaType and mediaTypeSize
	mediaType := []byte(strings.Split(file.Name(), ".")[1])
	header.MediaTypeSize = int8(len(mediaType))
	body.MediaType = mediaType

	// get payload and payloadSize
	fileBody := make([]byte, 0)
	_, err = file.Read(fileBody)
	if err != nil {
		return nil
	}
	header.PayloadSize = int64(len(fileBody))
	body.Payload = fileBody

	request := make([]byte, 0)
	return binaryRequest(request, header, body)
}

func binaryRequest(request []byte, header Header, body Body) []byte {
	return setFieldValue(request, header.htoByteSlice(), body.btoByteSlice())
}

func setFieldValue(base []byte, value ...interface{}) []byte {
	for _, v := range value {
		switch t := v.(type) {
		case int8:
			base = append(base, byte(t))
		case int16:
			base = append(base, byte(t))
		case int64:
			base = append(base, byte(t))
		case []byte:
			base = append(base, t...)
		}
	}

	return base
}

func (header *Header) htoByteSlice() []byte {
	byteSlice := make([]byte, 0)
	return setFieldValue(byteSlice, header.JsonSize, header.MediaTypeSize, header.PayloadSize)
}

func (body *Body) btoByteSlice() []byte {
	byteSlice := make([]byte, 0)
	return setFieldValue(byteSlice, body.Json, body.MediaType, body.Payload)
}
