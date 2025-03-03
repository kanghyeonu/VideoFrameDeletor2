package handler

import (
	"log"
	"math"
	"os"
	"strconv"
)

type videoHandler struct {
	h264ReadFileHandler  *h264ReadFileHandler
	h264WriteFileHandler *h264WriteFileHandler
	deleteOptions        *deleteOptions
	logFile              *os.File
	logger               *log.Logger
}

/*
"Parameters: 1 - 6
"  1. filename      string  : Input file name with .h264 extension
"  2. bytesToRemove int     : 0 to 100 if ratio is true, 0 to n if ratio is false
"  3. offset        int     : 1 to 100 Offset starting position for deletion in each Nalu
"  4. ratio         bool 	 : Ratio for processing (true: 1/false: 0)
"  5. reverse       bool    : Reverse the operation (true: 1/false: 0)
"  6. increment     int     : Increment value for offset
*/
// construct the videoHandler struct
func CreateVideoHandler(inputs []string) *videoHandler {
	filename := inputs[1]
	bytesToRemove := inputs[2]
	offset := inputs[3]
	ratio := inputs[4]
	reverse := inputs[5]
	increment := inputs[6]

	h264ReadFileHandler := createReadFileHandler(filename)
	deleteOptions := setDeleteOptions(bytesToRemove, offset, ratio, reverse, increment)
	return &videoHandler{
		h264ReadFileHandler: h264ReadFileHandler,
		deleteOptions:       deleteOptions,
	}
}

func (h *videoHandler) ResetFileHandler() {
	h.h264ReadFileHandler.reset()
	h.h264WriteFileHandler.init()
	h.logFile.Close()

}

func (h *videoHandler) SetWriteFileHandler(filename string) {
	h.h264WriteFileHandler = createWriteFileHandler(filename + ".h264")
}

func (h *videoHandler) GetDeleteOptions() (int, int, bool, bool, int) {
	return h.deleteOptions.bytesToRemove,
		h.deleteOptions.offset,
		h.deleteOptions.ratio,
		h.deleteOptions.reverse,
		h.deleteOptions.increment
}

func (h *videoHandler) CreateModifiedVideo(byteToRemove int, offset int, ratio bool, reverse bool, increment int) {

	numberOfNalu := 0 // number of NALU
	readSize := 0     // read size frome original file
	writeSize := 0    // write size to modified file
	deletedSize := 0
	maxNaluSize := math.Inf(-1) // max NALU size
	minNaluSize := math.Inf(1)  // set minNaluSize to the maximum value of int
	var naluLenSlice []int      // NALU slice

	// make channel for NALU
	nalu_chan := make(chan []byte)                 // channel for NALU
	go h.h264ReadFileHandler.getNalUnit(nalu_chan) // get NALU from original file

	for nalu := range nalu_chan {
		numberOfNalu++
		naluLen := len(nalu)
		if naluLen == 0 {
			//panic?
			continue
		}

		// update nalu info
		minNaluSize = math.Min(minNaluSize, float64(naluLen))
		maxNaluSize = math.Max(maxNaluSize, float64(naluLen))
		naluLenSlice = append(naluLenSlice, naluLen)
		readSize += naluLen

		// delete NALU & write file
		deletedNalu, delSize := deleteNaluByParams(nalu, byteToRemove, offset, ratio, reverse)
		ws := h.h264WriteFileHandler.writeNalUnit(deletedNalu)
		writeSize += ws
		deletedSize += delSize
	}

	// wirte info to log file
	h.logger.Printf("\nFile Path: %s\nOriginal File Name: %s\nNumber of NALU: %d\nRead size: %d\nWrite size: %d\nDeleted size: %d\nMax NALU size: %f\nMin NALU size: %f\nAverage Read NALU size: %f\nNumber of NALU in video: %d\n",
		"modified videos/"+strconv.Itoa(byteToRemove)+"_"+strconv.Itoa(offset)+"_"+strconv.FormatBool(ratio)+"_"+strconv.FormatBool(reverse)+"_"+strconv.Itoa(increment)+"/log.txt",
		h.h264ReadFileHandler.h264File.Name(),
		numberOfNalu,
		readSize,
		writeSize,
		deletedSize,
		maxNaluSize,
		minNaluSize,
		float64(readSize)/float64(numberOfNalu),
		len(naluLenSlice))

}

func (h *videoHandler) CreateLogFile(fileName string) *os.File {
	var file *os.File
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				log.Fatalln("Error creating log file: ", err)
			}
		}
	}
	h.logFile = file
	h.logger = log.New(h.logFile, "", log.Ldate|log.Ltime)
	return file
}

func (h *videoHandler) Close() {
	h.h264ReadFileHandler.close()
	h.h264WriteFileHandler.close()
	h.logFile.Close()
}
