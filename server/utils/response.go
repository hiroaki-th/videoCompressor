package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Header struct {
	Status        uint8
	JsonSize      uint16
	MediaTypeSize uint8
	PayloadSize   uint64
}

type Body struct {
	Json      []byte
	MediaType []byte
	Payload   []byte
}

type ResponseJson struct {
	Status   uint8
	FileName string
}

type ErrorJson struct {
	Status    uint8
	Message   string
	TimeStamp time.Time
}

func (header *Header) htoByteSlice() []byte {
	byteSlice := make([]byte, 0)

	return setFieldValue(byteSlice, header.JsonSize, header.MediaTypeSize, header.PayloadSize)
}

func (body *Body) btoByteSlice() []byte {
	byteSlice := make([]byte, 0)
	return setFieldValue(byteSlice, body.Json, body.MediaType, body.Payload)
}

func NewResponse(status uint8, file *os.File, errs ...error) []byte {

	if len(errs) > 0 {
		return ErrorResponse(status, errs...)
	}

	header := Header{}
	body := Body{}

	// get json and jsonSize
	filename := string([]byte(file.Name()))
	respJson := ResponseJson{
		Status:   status,
		FileName: filename,
	}
	byteJson, err := json.Marshal(respJson)
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

	response := make([]byte, 0)
	return binaryRequest(response, header, body)
}

func ErrorResponse(status uint8, errs ...error) []byte {

	header := Header{}
	body := Body{}

	// create error message
	message := ""
	for _, e := range errs {
		message = message + e.Error()
	}

	// create Json and get Size
	respJson := ErrorJson{
		Status:    status,
		Message:   message,
		TimeStamp: time.Now(),
	}
	byteJson, err := json.Marshal(respJson)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	header.JsonSize = uint16(len(byteJson))
	body.Json = byteJson

	// get mediaType and mediaTypeSize
	header.MediaTypeSize = uint8(0)
	body.MediaType = nil

	// get payload and payloadSize, limit is 10MB
	header.PayloadSize = uint64(0)
	body.Payload = nil

	response := make([]byte, 0)
	return binaryRequest(response, header, body)
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
		tmpResult := make([]byte, 8)
		binary.BigEndian.PutUint64(tmpResult, t)
		for len(tmpResult) != 8 {
			tmpResult = append([]byte{byte(0)}, tmpResult...)
		}
		result = tmpResult
	case uint32:
		tmpResult := make([]byte, 4)
		binary.BigEndian.PutUint32(tmpResult, t)
		for len(tmpResult) != 4 {
			tmpResult = append([]byte{byte(0)}, tmpResult...)
		}
		result = tmpResult
	case uint16:
		tmpResult := make([]byte, 2)
		binary.BigEndian.PutUint16(tmpResult, t)
		for len(tmpResult) != 2 {
			tmpResult = append([]byte{byte(0)}, tmpResult...)
		}
		result = tmpResult
	case uint8:
		result = append(result, t)
	}

	return result
}
