package handler

import (
	"fmt"
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

/*
 *@return []byte, int: deletedNalu, deletedBytes
 */
func deleteNaluByParams(nalu []byte, bytesToRemove int, offset int, ratio bool, reverse bool) ([]byte, int) {
	// nalu is empty
	if len(nalu) == 0 {
		fmt.Print("no data to be deleted")
		return nil, 0
	}

	// nothing to delete
	if bytesToRemove < 1 {
		fmt.Print("nalu copy mode: ", len(nalu), "Bytes are copied")
		return nalu, 0
	}

	// init startPatternLen
	var startPatternLen int
	if nalu[2] == start3BytePattern[2] {
		startPatternLen = 3
	} else {
		startPatternLen = 4
	}

	// pps and sps is excluded
	if nalu[startPatternLen]&0x1f == 7 || nalu[startPatternLen]&0x1f == 8 {
		return nalu, 0
	}

	//calculate offset position
	offsetPos := int(float64(offset*(len(nalu)))*0.01) + startPatternLen + 1 // +1 is nalu header size
	deleteSize := bytesToRemove

	var deletedNalu []byte
	copiedNalu := make([]byte, len(nalu))
	copy(copiedNalu, nalu)
	//  data ratio based, value must to be haved 1 ~ 99
	if ratio {
		// update deleteSize based on ratio
		deleteSize = int(float64(len(nalu)*bytesToRemove) * 0.01)
		if reverse { // delete from end of nalu
			deletedNalu = copiedNalu[:len(nalu)-deleteSize]
		} else {
			deletedNalu = copiedNalu[:offsetPos]
			// check if offset + deleteSize is over nalu's length
			if offsetPos+deleteSize > len(nalu) {
				return deletedNalu, len(nalu) - offsetPos
			}
			deletedNalu = append(deletedNalu, copiedNalu[offsetPos+deleteSize:]...)
		}

		// constant based, value don't be haved value over nalu's RBSP length
	} else {
		if len(nalu) < deleteSize {
			fmt.Print("Invalid value: constant must not be over nalu's RBSP length ", deleteSize)
			return nalu, 0
		}

		if reverse {
			deletedNalu = copiedNalu[:len(nalu)+1-deleteSize]
		} else {
			deletedNalu = copiedNalu[:offsetPos] // +1 is for includeing header
			if offsetPos+deleteSize > len(nalu) {
				return deletedNalu, len(nalu) - offsetPos
			}
			deletedNalu = append(deletedNalu, copiedNalu[offsetPos+deleteSize:]...)
		}
	}

	return deletedNalu, deleteSize
}
