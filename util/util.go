package util

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/**s 개수 확인 및 유효한지 검증하는 함수
 * @return string : 입력받은 파일명
 */
func ArgsParser() []string {
	args := os.Args
	if len(args) < 6 || len(args) > 6 {
		log.Fatalln("Usage: ./main.go {filename}.h264 {bytesToRemove} {offset} {ratio} {reverse}\n" +
			"Parameters:\n" +
			"  filename      string  : Input file name with .h264 extension\n" +
			"  bytesToRemove int     : 0 to 100 if ratio is true, 0 to n if ratio is false\n" +
			"  offset        int     : 5 to 100 Offset position for removal\n" +
			"  ratio         bool 	 : Ratio for processing (true: 1/false: 0)\n" +
			"  reverse       bool    : Reverse the operation (true: 1/false: 0)")
	}

	// validate the input parameters
	if isInvalidFileName(args[1]) ||
		isInvalidBytesToRemoveandRatio(args[2], args[4]) ||
		isInvalidOffset(args[3]) ||
		isInvalidReverse(args[5]) {

		log.Fatalln("Invalid Prameter. Please check the input format and try again.")
	}

	return args
}

func isInvalidFileName(filename string) bool {
	if strings.HasSuffix(filename, ".h264") {
		return false
	}
	fmt.Println("Please check the file foramt value and try again.")
	return true
}
func isInvalidBytesToRemove(bytesToRemove string) bool {
	bytesToRemoveInt, err := strconv.Atoi(bytesToRemove)
	if err != nil {
		fmt.Println("Please check the bytesToRemove value and try again.\nError: ", err)
		return true
	}
	if bytesToRemoveInt < 0 {
		fmt.Println("bytesToRemove must be a non-negative integer.")
		return true
	}
	return false
}

func isInvalidBytesToRemoveandRatio(bytesToRemove string, ratio string) bool {
	bytesToRemoveInt, err := strconv.Atoi(bytesToRemove)
	if err != nil {
		fmt.Println("Please check the bytesToRemove value and try again.\nError: ", err)
		return true
	}
	ratioInt, err := strconv.Atoi(ratio)
	if err != nil {
		fmt.Println("Please check the ratio value and try again.\nError: ", err)
		return true
	}

	if ratioInt == 1 {
		if bytesToRemoveInt < 0 || bytesToRemoveInt > 100 {
			fmt.Println("Please check the bytesToRemove value and try again.\n" +
				"ratio is true, but bytesToRemove is not in the range of 0 to 100.")
			return true
		}
	} else {
		if bytesToRemoveInt < 0 {
			fmt.Println("Please check the bytesToRemove value and try again.\n" +
				"ratio is false, but bytesToRemove is not in the range of 0 to n.")
			return true
		}
	}
	return false
}

func isInvalidOffset(offset string) bool {
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		fmt.Println("Please check the offset value and try again.\nError: ", err)
		return true
	}
	if offsetInt < 5 || offsetInt > 100 {
		fmt.Println("Offset must be in the range of 5 to 100.")
		return true
	}
	return false
}

func isInvalidReverse(reverse string) bool {
	reverseInt, err := strconv.Atoi(reverse)
	if err != nil {
		fmt.Println("Please check the reverse value and try again.\nError: ", err)
		return true
	}
	if reverseInt != 0 && reverseInt != 1 {
		fmt.Println("Reverse must be 0 (false) or 1 (true).")
		return true
	}
	return false
}
