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
	// TODO
	if len(nalu) == 0 {
		fmt.Print("no data to be deleted")
		return nil, 0
	} else if p.bytesToRemove < 1 {
		fmt.Print("nalu copy mode: ", len(nalu), "Bytes are copied")
		return nalu, 0
	}
	N := p.bytesToRemove

	var startPatternLen int
	if nalu[2] == start3BytePattern[2] {
		startPatternLen = 3
	} else {
		startPatternLen = 4
	}

	// pps and sps is exception
	if nalu[startPatternLen]&0x1f == 7 || nalu[startPatternLen]&0x1f == 8 {
		return nalu, 0
	}

	//TODO
	offset := int(float64(offset*(len(nalu)-startPatternLen-1))*0.01) + startPatternLen + 1 // +1 is nalu header size
	ratio := ratio
	reverse := reverse

	var deletedNalu []byte
	copiedNalu := make([]byte, len(nalu))
	copy(copiedNalu, nalu)
	//  data ratio based, value must to be haved 1 ~ 99
	if ratio {
		if N > 99 || N < 1 {
			fmt.Print("Invalid value: ratio must be 1~99 ", N)
			return nalu
		}

		sizeToDelete := int(float64(len(nalu)*N) * 0.01)
		if !reverse {
			deletedNalu = copiedNalu[:offset]
			if offset+sizeToDelete > len(nalu) {
				return deletedNalu
			}

			deletedNalu = append(deletedNalu, copiedNalu[offset+sizeToDelete:]...)
		} else {
			deletedNalu = copiedNalu[:len(nalu)-sizeToDelete+1]
		}

		// constant based, value don't be haved value over nalu's RBSP length
	} else {
		if len(nalu)-1 < N {
			fmt.Print("Invalid value: constant must not be over nalu's RBSP length ", N)
			return nalu
		}

		if !reverse {
			deletedNalu = copiedNalu[:offset] // +1 is for includeing header
			if offset+N > len(nalu)-1 {
				return deletedNalu
			}
			deletedNalu = append(deletedNalu, copiedNalu[offset+N:]...)

		} else {
			deletedNalu = copiedNalu[:len(nalu)+1-N]
		}
	}

	return deletedNalu
}
