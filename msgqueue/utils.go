package msgqueue

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetFilePathAndFileInmyposchainFromDir(dir string, maxFileSize int) (filePath string, fileInmyposchain int, err error) {
	if _, fileInmyposchain, err = getFilePathAndInmyposchain(dir, maxFileSize); err != nil {
		return
	}

	fileSize := GetFileSize(GetFileName(dir, fileInmyposchain))
	if fileSize < maxFileSize {
		return GetFileName(dir, fileInmyposchain), fileInmyposchain, nil
	}
	return GetFileName(dir, fileInmyposchain+1), fileInmyposchain + 1, nil
}

func getFilePathAndInmyposchain(dir string, height int) (filePath string, fileInmyposchain int, err error) {
	fileNames, err := getAllFilesFromDir(dir)
	if err != nil {
		return "", -1, err
	}
	if len(fileNames) == 0 {
		return GetFileName(dir, 0), 0, nil
	}
	fileInmyposchain = getMaxInmyposchainFromFiles(fileNames)
	return GetFileName(dir, fileInmyposchain), fileInmyposchain, nil
}

func getAllFilesFromDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return nil, err
		}
		return nil, nil
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if strings.HasPrefix(name, filePrefix) {
				fileNames = append(fileNames, name)
			}
		}
	}
	return fileNames, nil
}

func getMaxInmyposchainFromFiles(fileNames []string) int {
	fileInmyposchain := 0
	for _, fileName := range fileNames {
		vals := strings.Split(fileName, "-")
		if len(vals) == 2 {
			if inmyposchain, err := strconv.Atoi(vals[1]); err == nil {
				if inmyposchain > fileInmyposchain {
					fileInmyposchain = inmyposchain
				}
			}
		}
	}
	return fileInmyposchain
}

func getMinInmyposchainFromFiles(fileNames []string) int {
	fileInmyposchain := math.MaxInt64
	for _, fileName := range fileNames {
		vals := strings.Split(fileName, "-")
		if len(vals) == 2 {
			if inmyposchain, err := strconv.Atoi(vals[1]); err == nil {
				if inmyposchain < fileInmyposchain {
					fileInmyposchain = inmyposchain
				}
			}
		}
	}
	return fileInmyposchain
}

func GetFileName(dir string, fileInmyposchain int) string {
	return dir + "/" + filePrefix + strconv.Itoa(fileInmyposchain)
}

func GetFileSize(filePath string) int {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) || info.IsDir() {
		return -1
	}

	return int(info.Size())
}

func GetFileLeastHeightInDir(dir string) (string, int64, error) {
	files, err := getAllFilesFromDir(dir)
	if err != nil {
		return "", -1, err
	}
	inmyposchain := getMinInmyposchainFromFiles(files)
	in, err := openFile(GetFileName(dir, inmyposchain))
	if err != nil {
		return "", -1, err
	}
	defer in.Close()

	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil {
		return "", -1, err
	}
	return GetFileName(dir, inmyposchain), getHeight(line), nil
}

func getHeight(data string) int64 {
	if len(data) == 0 {
		return -1
	}
	vals := strings.Split(data, "#")
	if vals[0] != "height_info" {
		panic("The first line of data in the file should be [height_info] msg")
	}
	var info NewHeightInfo
	err := json.Unmarshal([]byte(vals[1]), &info)
	if err != nil {
		panic(fmt.Sprintf("json unmarshal error : %s\n", err.Error()))
	}
	return info.Height
}

func openFile(filePath string) (*os.File, error) {
	if s, err := os.Stat(filePath); os.IsNotExist(err) {
		return os.Create(filePath)
	} else if s.IsDir() {
		return nil, fmt.Errorf("Need to give the file path ")
	} else {
		return os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0666)
	}
}

func FillMsgs(ctx sdk.Context, key string, msg interface{}) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(EventTypeMsgQueue, sdk.NewAttribute(key, string(bytes))))
}
