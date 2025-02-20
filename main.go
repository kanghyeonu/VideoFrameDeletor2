package main

import (
	"fmt"
	"os"
	"strconv"
	"videoframedetector2/handler"
	"videoframedetector2/util"
)

func main() {
	start()
}

func start() {
	// get file name from command line and validate it
	inputs, err := util.ArgsParser([]string{})
	if err != nil {
		print(err)
		os.Exit(1)
	}
	// create file handler using the validated inputs
	h := handler.CreateVideoHandler(inputs)
	print("success to create video handler")

	byteToRemove, offset, ratio, reverse, increment := h.GetDeleteOptions()

	dirName := "modified videos/" + strconv.Itoa(byteToRemove) + "_" + strconv.Itoa(offset) + "_" + strconv.FormatBool(ratio) + "_" + strconv.FormatBool(reverse) + "_" + strconv.Itoa(increment)
	err = util.CreateDirectory(dirName)
	if err != nil {
		os.Exit(1)
	}

	// create modified videos
	for start_offset := offset; start_offset <= 100; start_offset += increment {
		// create modified video name
		// modified video name format: "{offset}.h264"
		// increment the offset by the increment value
		// e.g. increment = 5, offset = 5
		// 		"5.h264", "10.h264", "15.h264", ..., "100.h264"
		modifiedVideoName := dirName + "/" + strconv.Itoa(start_offset) + ".h264"
		fmt.Print(modifiedVideoName + " processing...\n")

		// set
		h.SetWriteFileHandler(modifiedVideoName)
		h.CreateModifiedVideo(h.GetDeleteOptions()) // TODO
	}

	//
}
