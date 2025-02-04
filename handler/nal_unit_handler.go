package handler

import (
	"strconv"
)

type deleteOptions struct {
	bytesToRemove int
	offset        int
	ratio         bool
	reverse       bool
}

func setDeleteOptions(bytesToRemove string, offset string, ratio string, reverse string) *deleteOptions {
	bytesToRemoveInt, _ := strconv.Atoi(bytesToRemove)

	offsetInt, _ := strconv.Atoi(offset)

	ratioInt, _ := strconv.Atoi(ratio)
	ratiobool := ratioInt == 1

	reverseInt, _ := strconv.Atoi(reverse)
	reversebool := reverseInt == 1

	return &deleteOptions{
		bytesToRemove: bytesToRemoveInt,
		offset:        offsetInt,
		ratio:         ratiobool,
		reverse:       reversebool,
	}
}
