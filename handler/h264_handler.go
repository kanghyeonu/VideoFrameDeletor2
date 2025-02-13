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

// openFile opens the file with the given file name
func openFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Please check the file name and try again.\nError: ", err)
		return nil
	}
	return file
}

// createFile creates a file with the given file name for modified video
func createFile(fileName string) *os.File {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("Please check the file name and try again.\nError: ", err)
		return nil
	}
	return file
}

// createReadFileHandler creates a file handler for reading the h264 file
func createReadFileHandler(fileName string) *h264ReadFileHandler {
	file := openFile("original videos/" + fileName)
	return &h264ReadFileHandler{
		h264File:         file,
		fileReader:       bufio.NewReader(file),
		fileReaderBuffer: make([]byte, 4096),
	}
}

// createWriteFileHandler creates a file handler for writing the h264 file
func createWriteFileHandler(fileName string) *h264WriteFileHandler {
	file := createFile("modified videos/" + fileName)
	writer := bufio.NewWriter(file)
	return &h264WriteFileHandler{
		objectFile: file,
		fileWriter: writer,
	}
}
