package api

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const path string = "./tmp/"

type FileJson struct {
	Name       string `json:"name"`
	Extension  string `json:"extension"`
	Resolution string `json:"resolution"`
	FromSecond int    `json:"from-second"`
	ToSecond   int    `json:"to-second"`
	VF         string `json:"vf"`
}

func SaveFile(buff []byte) (*os.File, *FileJson, error) {
	jsonSize := int(binary.BigEndian.Uint16(buff[:2]))
	mediaTypeSize := int(int8(buff[2]))
	payloadSize := int(binary.BigEndian.Uint64(buff[3:11]))

	jsonBin := buff[11 : 11+jsonSize]
	_ = buff[11+jsonSize : 11+jsonSize+mediaTypeSize]
	payload := buff[11+jsonSize+mediaTypeSize : 11+jsonSize+mediaTypeSize+payloadSize]

	jsonData := FileJson{}
	err := json.Unmarshal(jsonBin, &jsonData)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.Create(path + jsonData.Name)
	if err != nil {
		return nil, nil, err
	}

	_, err = file.Write(payload)
	if err != nil {
		return nil, nil, err
	}

	return file, &jsonData, nil
}

func FormatFile(file *os.File, fileJson *FileJson) (*os.File, error) {

	input := path + fileJson.Name
	output := path + strings.Split(fileJson.Name, ".")[0] + fileJson.Extension

	fileCmd := fmt.Sprintf(`ffmpeg -i %s`, input)

	if fileJson.Resolution != "auto" {
		fileCmd = fileCmd + fmt.Sprintf(` -s %s`, fileJson.Resolution)
	}

	if fileJson.VF != "auto" {
		fileCmd = fileCmd + fmt.Sprintf(` -vf %s`, fileJson.VF)
	}

	if fileJson.FromSecond != 0 {
		fileCmd = fileCmd + fmt.Sprintf(` -ss %d -to  %d`, fileJson.FromSecond, fileJson.ToSecond)
	}

	fileCmd = fileCmd + fmt.Sprintf(` %s`, output)

	_, err := exec.Command("sh", "-c", fileCmd).CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 1 {
			fmt.Println(exec.Command("sh", "-c", fileCmd).String())
			return nil, err
		}
	}

	err = os.Remove(path + fileJson.Name)
	if err != nil {
		return nil, err
	}

	formattedFile, err := os.Open(output)
	if err != nil {
		return nil, err
	}

	return formattedFile, nil
}

func getTotalSize(buff []byte) int {
	jsonSize := int(binary.BigEndian.Uint16(buff[:2]))
	mediaTypeSize := int(int8(buff[2]))
	payloadSize := int(binary.BigEndian.Uint64(buff[3:11]))

	return 11 + jsonSize + mediaTypeSize + payloadSize
}
