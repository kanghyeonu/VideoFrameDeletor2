package main

import (
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
		return
	}
	// create file handler using the validated inputs
	handler := handler.CreateFileHandler(inputs)
	print("success to create file handler ", handler)

	// TODO
}
