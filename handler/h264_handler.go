package handler

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
)

// variables for KMP algorithm
// pattern start sequence of NALU
var (
	init4Bytetable    = []int{0, 1, 2, 0}
	init3Bytetable    = []int{0, 1, 0}
	start4BytePattern = []byte{0, 0, 0, 1}
	start3BytePattern = []byte{0, 0, 1}
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

// reset resets the file handler to the beginning of the file
func (h *h264ReadFileHandler) reset() {
	_, err := h.h264File.Seek(0, 0)
	if err != nil {
		log.Fatalln(err)
	}
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
	file := createFile(fileName)
	writer := bufio.NewWriter(file)
	return &h264WriteFileHandler{
		objectFile: file,
		fileWriter: writer,
	}
}

func (h *h264ReadFileHandler) close() {
	h.h264File.Close()
}

func (h *h264WriteFileHandler) close() {
	h.objectFile.Close()
}

func (h *h264WriteFileHandler) init() {
	h.fileWriter.Flush()
	h.objectFile.Close()

	h.objectFile = nil
	h.fileWriter = nil
}

func (h *h264WriteFileHandler) writeNalUnit(nalu []byte) int {
	ws, err := h.fileWriter.Write(nalu)
	if err != nil {
		log.Fatalln(err)
	}
	return ws
}

func (h *h264ReadFileHandler) getNalUnit(nalu_chan chan []byte) {
	defer close(nalu_chan)

	h.remainingBytes = []byte{}
	total := 0
	for {
		// read file
		readDataBytes, err := io.ReadFull(h.fileReader, h.fileReaderBuffer[:cap(h.fileReaderBuffer)])
		h.fileReaderBuffer = h.fileReaderBuffer[:readDataBytes] // stream data, not nalu yet
		if err != nil {
			if err == io.EOF {
				break
			}
			if err != io.ErrUnexpectedEOF {
				log.Fatalln(err)
				panic(err)
			}
		}

		total += readDataBytes
		currentLength := len(h.remainingBytes)
		h.remainingBytes = append(h.remainingBytes, h.fileReaderBuffer...)

		var startPositions = []int{}

		// 버퍼에 남은 데이터가 3byte 이상이면 -> 이어붙인 데이터의 시작 위치가 0th 인덱스일 경우 고려
		if currentLength > 3 {
			offset := currentLength - 3
			startPositions = findStartSequencePosition(h.remainingBytes[offset:])
			for idx, val := range startPositions {
				startPositions[idx] = val + offset
			}

			// 3byte 미만이면 무조건
		} else {
			startPositions = findStartSequencePosition(h.remainingBytes)
			if len(startPositions) > 2 {
				startPositions = startPositions[1:]
			}
		}

		if len(startPositions) >= 1 && startPositions[0] != 0 {
			//start := startPositions[0]
			start := 0
			for _, end := range startPositions {
				if start-end == 0 {

				}
				nalu_chan <- h.remainingBytes[start:end]
				start = end
			}
			//pf("[%d-]: %x\n",start ,nalu[:4])
			h.remainingBytes = h.remainingBytes[start:]
			//pf("[%d-]: %x\n",start , h.remainingBytes[:4])
			//h.remainingBytes = append(h.remainingBytes, h.fileReaderBuffer[start:]...)
		}

	}
	//last nal unit
	nalu_chan <- h.remainingBytes
}

/*
 * find start sequence position in data
 * @return []byte, int: deletedNalu, deletedBytes
 */
func findStartSequencePosition(data []byte) []int {
	ch := make(chan int, 2)
	go kmp(data, start4BytePattern, ch)
	go kmp(data, start3BytePattern, ch)

	pos := []int{}
	count := 0
	for offset := range ch {
		//case: go routines is end
		if offset == -1 {
			count++
			if count == 2 {
				break
			}
		} else {
			pos = append(pos, offset)
		}
	}
	close(ch)

	// no start point
	if len(pos) == 0 {
		return pos
	}

	sort.Ints(pos)

	startPoints := []int{}
	for i := 0; i < len(pos)-1; i++ {
		// case: 4byte pattern
		// 4 byte start sequence is includes a 3 byte start sequence
		if pos[i+1]-pos[i] == 1 {
			startPoints = append(startPoints, pos[i])

			// case: 3byte pattern
		} else {
			startPoints = append(startPoints, pos[i+1])
		}
	}

	// make to unique slice
	keys := make(map[int]struct{})
	result := make([]int, 0)
	for _, val := range startPoints {
		if _, ok := keys[val]; ok {
			continue
		} else {
			keys[val] = struct{}{}
			result = append(result, val)
		}
	}
	// no need last nal unit start point
	// becase cannot know end of the last nalu
	return result
}

func kmp(data []byte, pattern []byte, ch chan int) {
	byteLength := len(data)
	patternLength := len(pattern)

	var table []int

	if patternLength == 4 {
		table = init4Bytetable
	} else {
		table = init3Bytetable
	}

	i := 0
	for i < byteLength {
		j := 0
		for i < byteLength && j < patternLength {
			if data[i] == pattern[j] {
				i++
				j++
			} else {
				if j != 0 {
					j = table[j-1]
				} else {
					i++
				}
			}
		}
		// case: pattern pattern
		if j == patternLength {
			ch <- i - j
		}

		// case: end of data
		if i == byteLength {
			ch <- -1
			return
		}
	}
}
