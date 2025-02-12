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
	inputs := util.ArgsParser()
	handler := handler.CreateFileHandler(inputs)
	print(handler)

	// TODO
}
