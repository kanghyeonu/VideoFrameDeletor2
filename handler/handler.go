package handler

type fileHandler struct {
	h264ReadFileHandler  *h264ReadFileHandler
	h264WriteFileHandler *h264WriteFileHandler
	deleteOptions        *deleteOptions
}

/*
"Parameters: 1 - 5
"  1. filename      string  : Input file name with .h264 extension
"  2. bytesToRemove int     : 0 to 100 if ratio is true, 0 to n if ratio is false
"  3. offset        int     : 5 to 100 Offset position for removal
"  4. ratio         bool 	 : Ratio for processing (true: 1/false: 0)
"  5. reverse       bool    : Reverse the operation (true: 1/false: 0)
*/
// construct the fileHandler struct
func CreateFileHandler(inputs []string) *fileHandler {
	filename := inputs[1]
	bytesToRemove := inputs[2]
	offset := inputs[3]
	ratio := inputs[4]
	reverse := inputs[5]

	h264ReadFileHandler := createReadFileHandler(filename)
	h264WriteFileHandler := createWriteFileHandler(filename)
	deleteOptions := setDeleteOptions(bytesToRemove, offset, ratio, reverse)
	return &fileHandler{
		h264ReadFileHandler:  h264ReadFileHandler,
		h264WriteFileHandler: h264WriteFileHandler,
		deleteOptions:        deleteOptions,
	}
}

func (fh *fileHandler) CreateModifiedVideo() {
	// init parameters
	numberOfNalu := 0                 // number of NALU
	readSize := 0                     // read size frome original file
	writeSize := 0                    // write size to modified file
	maxNaluSize := 0                  // max NALU size
	minNaluSize := int(^uint(0) >> 1) // set minNaluSize to the maximum value of int

	// make channel for NALU
	nalu_chan := make(chan []byte) // channel for NALU
	go fh.h264ReadFileHandler.getNalUnit(nalu_chan)
	for nalu := range nalu_chan {
		// TDO process NALU
	}
}
