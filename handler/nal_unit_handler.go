package handler

import (
	"strconv"
)

type deleteOptions struct {
	bytesToRemove int
	offset        int
	ratio         bool
	reverse       bool
	increment     int
}

func setDeleteOptions(bytesToRemove string, offset string, ratio string, reverse string, increment string) *deleteOptions {
	bytesToRemoveInt, _ := strconv.Atoi(bytesToRemove)
	offsetInt, _ := strconv.Atoi(offset)
	ratioInt, _ := strconv.Atoi(ratio)
	ratioBool := ratioInt == 1
	reverseInt, _ := strconv.Atoi(reverse)
	reverseBool := reverseInt == 1
	incrementInt, _ := strconv.Atoi(increment)

	return &deleteOptions{
		bytesToRemove: bytesToRemoveInt,
		offset:        offsetInt,
		ratio:         ratioBool,
		reverse:       reverseBool,
		increment:     incrementInt,
	}
}
