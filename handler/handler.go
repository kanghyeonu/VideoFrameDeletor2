package handler

import (
	"fmt"
	"math"
)

type videoHandler struct {
	h264ReadFileHandler  *h264ReadFileHandler
	h264WriteFileHandler *h264WriteFileHandler
	deleteOptions        *deleteOptions
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
	fh.h264WriteFileHandler = createWriteFileHandler(filename)
}

func (fh *videoHandler) GetDeleteOptions() (int, int, bool, bool, int) {
	return fh.deleteOptions.bytesToRemove,
		fh.deleteOptions.offset,
		fh.deleteOptions.ratio,
		fh.deleteOptions.reverse,
		fh.deleteOptions.increment
}

func (fh *videoHandler) CreateModifiedVideo(byteToRemove int, offset int, ratio bool, reverse bool) {

	numberOfNalu := 0           // number of NALU
	readSize := 0               // read size frome original file
	writeSize := 0              // write size to modified file
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

		// delete NALU
		deletedNalu, deletedSize := deleteNaluByParams(nalu, byteToRemove, offset, ratio, reverse)

		// TODO : write NALU to modified file

	}
	close(nalu_chan)

	// print nalu info
	// 나중에 리스트 형식으로 반환?
	fmt.Println("Number of NALU: ", numberOfNalu)
	fmt.Println("Read size: ", readSize)
	fmt.Println("Write size: ", writeSize)
	fmt.Print("Deleted size: ", readSize-writeSize)
	fmt.Println("Max NALU size: ", maxNaluSize)
	fmt.Println("Min NALU size: ", minNaluSize)
	fmt.Println("number of NALU in video: ", len(naluLenSlice))

}
