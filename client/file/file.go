package file

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type FileJson struct {
	Name       string `json:"name"`
	Extension  string `json:"extension"`
	Resolution string `json:"resolution"`
	FromSecond int    `json:"from-second"`
	ToSecond   int    `json:"to-second"`
	VF         string `json:"vf"`
}

func SelectFile(reader *bufio.Reader) (*os.File, error) {
	fmt.Printf("please input filename you want to change format\n")
	filename, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	filename = strings.Trim(filename, "\n")
	fmt.Printf("searching %s...\n\n", filename)

	err = FindFiles(filename)
	if err != nil {
		return nil, err
	}

	fmt.Printf("select and input file path from above filepath list\n")
	fullFilePath, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fullFilePath = strings.Trim(fullFilePath, "\n")

	file, err := os.Open(fullFilePath)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n")
	return file, nil
}

func FindFiles(filename string) error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	findCmd := fmt.Sprintf(`find %s -type f -name %q -not -path "*/.*/*" 2>/dev/null`, home, filename)
	result, err := exec.Command("sh", "-c", findCmd).CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 1 {
			fmt.Println(exec.Command("sh", "-c", findCmd).String())
			return err
		}
	}

	fmt.Println(string(result))
	return nil
}

func SelectFormat(reader *bufio.Reader, file *os.File) (*FileJson, error) {
	fileJson := FileJson{}

	filepath := string([]byte(file.Name()))
	filename := strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
	fileJson.Name = filename

	extension, err := selection(reader, "how do you want format file? tell me [.extension]")
	if err != nil {
		return nil, err
	}
	fileJson.Extension = extension

	resolution, err := selection(reader, "tell me resolution you want. please input like [width(int)xheight[int]]. if you want as it is, please input 'auto'")
	if err != nil {
		return nil, err
	}
	fileJson.Resolution = resolution

	vf, err := selection(reader, "tell me aspect ration you want. please input like [width(int):height[int]]. if you want auto size width/height, input like [width:-1] which means height is auto sizing.\nif you want as it is, please input 'auto'")
	if err != nil {
		return nil, err
	}
	fileJson.VF = vf

	isCut, err := selection(reader, "do you want cut video in specific moment? [y/n]")
	if err != nil {
		return nil, err
	}

	if isCut[0] == 'y' {
		from, err := selection(reader, "where do you want to start video? tell me in second")
		if err != nil {
			return nil, err
		}
		fromSecond, err := strconv.Atoi(from)
		if err != nil {
			return nil, err
		}
		fileJson.FromSecond = fromSecond

		to, err := selection(reader, "where do you want to finish video? tell me in second\n")
		if err != nil {
			return nil, err
		}
		toSecond, err := strconv.Atoi(to)
		if err != nil {
			return nil, err
		}
		fileJson.ToSecond = toSecond
	}

	if isValid := ValidJson(&fileJson); !isValid {
		fmt.Printf("please input format you like one more time\n\n")
		return SelectFormat(reader, file)
	}

	if isOk, _ := ConfirmFileJson(reader, &fileJson); !isOk {
		fmt.Printf("please input format you like one more time\n\n")
		return SelectFormat(reader, file)
	}

	return &fileJson, nil
}

func selection(reader *bufio.Reader, question string) (string, error) {
	fmt.Printf("%s\n", question)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if answer == "\n" {
		return selection(reader, question)
	}
	answer = strings.Trim(answer, "\n")
	fmt.Printf("\n")
	return answer, nil
}

func ValidJson(fileJson *FileJson) bool {
	if fileJson.Extension == "" || fileJson.Extension[0] != '.' {
		fmt.Println("input extension. do not forget '.' before extension")
		return false
	}

	if fileJson.Resolution != "auto" && (!validIntStr(strings.Split(fileJson.Resolution, "x")...) || len(strings.Split(fileJson.Resolution, "x")) != 2) {
		fmt.Printf("Error: input resolution in right format, your input was '%s'\n", fileJson.Resolution)
		return false
	}

	if fileJson.VF != "auto" && (!validIntStr(strings.Split(fileJson.VF, ":")...) || len(strings.Split(fileJson.VF, ":")) != 2) {
		fmt.Printf("Error: input aspect ration in right format, your input was '%s'\n", fileJson.VF)
		return false
	}

	return true
}

func validIntStr(intStr ...string) bool {

	for _, str := range intStr {
		if _, err := strconv.Atoi(str); err != nil {
			return false
		}
	}

	return true
}

func ConfirmFileJson(reader *bufio.Reader, fileJson *FileJson) (bool, error) {
	fmt.Printf("Do you confirm to send this request? [y/n]\n")
	fmt.Printf("Name: %s\n", fileJson.Name)
	fmt.Printf("Extension: %s\n", fileJson.Extension)
	fmt.Printf("Resolution: %s\n", fileJson.Resolution)
	fmt.Printf("VF: %s\n", fileJson.VF)

	if fileJson.ToSecond != 0 {
		fmt.Printf("FromSecond: %d\n", fileJson.FromSecond)
		fmt.Printf("ToSecond: %d\n", fileJson.ToSecond)
	}

	answer, err := reader.ReadString('\n')
	if err != nil || answer[0] == 'n' {
		return false, err
	}

	return true, nil
}
