package handler

import (
	"log"
	"math"
	"os"
)

type videoHandler struct {
	h264ReadFileHandler  *h264ReadFileHandler
	h264WriteFileHandler *h264WriteFileHandler
	deleteOptions        *deleteOptions
	logFile              *os.File
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

func (fh *videoHandler) SetWriteFileHandler(filename string) {
	fh.h264WriteFileHandler = createWriteFileHandler(filename + ".264")
	fh.logFile = createLogFile(filename + ".txt")

}

func (fh *videoHandler) GetDeleteOptions() (int, int, bool, bool, int) {
	return fh.deleteOptions.bytesToRemove,
		fh.deleteOptions.offset,
		fh.deleteOptions.ratio,
		fh.deleteOptions.reverse,
		fh.deleteOptions.increment
}

func (fh *videoHandler) CreateModifiedVideo(byteToRemove int, offset int, ratio bool, reverse bool) {

	numberOfNalu := 0 // number of NALU
	readSize := 0     // read size frome original file
	writeSize := 0    // write size to modified file
	deletedSize := 0
	maxNaluSize := math.Inf(-1) // max NALU size
	minNaluSize := math.Inf(1)  // set minNaluSize to the maximum value of int
	var naluLenSlice []int      // NALU slice

	// make channel for NALU
	nalu_chan := make(chan []byte)                  // channel for NALU
	go fh.h264ReadFileHandler.getNalUnit(nalu_chan) // get NALU from original file

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
		ws := fh.h264WriteFileHandler.writeNalUnit(deletedNalu)
		writeSize += ws
		deletedSize += delSize
	}
	close(nalu_chan)

	// print nalu info
	// 나중에 리스트 형식으로 반환?
	// log파일 생성? csv?
	logger := log.New(fh.logFile, "", log.Ldate|log.Ltime)
	logger.Println("File name: ", fh.h264WriteFileHandler.objectFile.Name())
	logger.Println("Number of NALU: ", numberOfNalu)
	logger.Println("Read size: ", readSize)
	logger.Println("Write size: ", writeSize)
	logger.Println("Deleted size: ", deletedSize)
	logger.Println("Max NALU size: ", maxNaluSize)
	logger.Println("Min NALU size: ", minNaluSize)
	logger.Println("number of NALU in video: ", len(naluLenSlice))

}

func createLogFile(fileName string) *os.File {
	var file *os.File
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(fileName)
			if err != nil {
				log.Fatalln("Error creating log file: ", err)
			}
		}
	}

	return file
}
