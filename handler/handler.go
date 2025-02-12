package handler

type fileHandler struct {
	h264ReadFileHandler  *h264ReadFileHandler
	h264WriteFileHandler *h264WriteFileHandler
	deleteOptions        *deleteOptions
}

// construct the fileHandler struct
func CreateFileHandler(inputs []string) *fileHandler {
	h264ReadFileHandler := createReadFileHandler(inputs[1])
	h264WriteFileHandler := createWriteFileHandler(inputs[1]) // TODO
	deleteOptions := setDeleteOptions(inputs[2], inputs[3], inputs[4], inputs[5])
	return &fileHandler{
		h264ReadFileHandler:  h264ReadFileHandler,
		h264WriteFileHandler: h264WriteFileHandler,
		deleteOptions:        deleteOptions,
	}
}
