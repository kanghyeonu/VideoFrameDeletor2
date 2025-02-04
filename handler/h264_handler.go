package handler

import (
	"bufio"
	"log"
	"os"
)

type h264ReadFileHandler struct {
	h264File         *os.File
	fileReader       *bufio.Reader
	fileReaderBuffer []byte
	remainingBytes   []byte
}

type h264WriteFileHandler struct {
	objectFile *os.File
	fileWriter *bufio.Writer
}

func openFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Please check the file name and try again.\nError: ", err)
		return nil
	}
	return file
}

func createReadFileHandler(fileName string) *h264ReadFileHandler {
	file := openFile("video file/" + fileName)
	return &h264ReadFileHandler{
		h264File:         file,
		fileReader:       bufio.NewReader(file),
		fileReaderBuffer: make([]byte, 4096),
	}
}

func createWriteFileHandler(fileName string) *h264WriteFileHandler {
	//TODO: create objectfile according to the parameter

	return nil
}
