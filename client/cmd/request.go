package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	myFile "videoCompressorClient/file"
)

type Header struct {
	JsonSize      uint16
	MediaTypeSize uint8
	PayloadSize   uint64
}

type Body struct {
	Json      []byte
	MediaType []byte
	Payload   []byte
}

func (header *Header) htoByteSlice() []byte {
	byteSlice := make([]byte, 0)

	return setFieldValue(byteSlice, header.JsonSize, header.MediaTypeSize, header.PayloadSize)
}

func (body *Body) btoByteSlice() []byte {
	byteSlice := make([]byte, 0)
	return setFieldValue(byteSlice, body.Json, body.MediaType, body.Payload)
}

func CreateRequest(file *os.File) []byte {

	header := Header{}
	body := Body{}

	// get json and jsonSize
	filename := string([]byte(file.Name()))
	fileJson := myFile.FileJson{
		Name: filename,
	}
	byteJson, err := json.Marshal(fileJson)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	header.JsonSize = uint16(len(byteJson))
	body.Json = byteJson

	// get mediaType and mediaTypeSize
	mediaType := []byte(strings.Split(file.Name(), ".")[1])
	header.MediaTypeSize = uint8(len(mediaType))
	body.MediaType = mediaType

	// get payload and payloadSize, limit is 10MB
	fileBody := make([]byte, 104857600)
	size, err := file.Read(fileBody)
	if err != nil {
		return nil
	}
	header.PayloadSize = uint64(size)
	body.Payload = fileBody[:size]

	request := make([]byte, 0)
	return binaryRequest(request, header, body)
}

func binaryRequest(request []byte, header Header, body Body) []byte {
	return setFieldValue(request, header.htoByteSlice(), body.btoByteSlice())
}

func setFieldValue(base []byte, value ...interface{}) []byte {
	for _, v := range value {
		switch t := v.(type) {
		case uint8:
			base = append(base, t)
		case uint16:
			base = append(base, getByteSliceFromNumber(t)...)
		case uint64:
			base = append(base, getByteSliceFromNumber(t)...)
		case []byte:
			base = append(base, t...)
		}
	}
	return base
}

func getByteSliceFromNumber(number interface{}) []byte {
	result := make([]byte, 0)
	switch t := number.(type) {
	case uint64:
		for t > 255 {
			result = append([]byte{byte(t)}, result...)
			t = t / 255
		}
		result = append([]byte{byte(t)}, result...)
		for len(result) != 8 {
			result = append([]byte{byte(0)}, result...)
		}
	case uint32:
		for t > 255 {
			result = append([]byte{byte(t)}, result...)
			t = t / 255
		}
		result = append([]byte{byte(t)}, result...)
		for len(result) != 4 {
			result = append([]byte{byte(0)}, result...)
		}
	case uint16:
		for t > 255 {
			result = append([]byte{byte(t)}, result...)
			t = t / 255
		}
		result = append([]byte{byte(t)}, result...)
		for len(result) != 2 {
			result = append([]byte{byte(0)}, result...)
		}
	case uint8:
		result = append(result, t)
	}

	return result
}
