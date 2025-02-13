package util

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
 * 파라미터 검증 함수수
 * @return []string : 입력받은 파라미터 문자열 배열열
 */
type argsError struct {
	Message string
}

func (e *argsError) Error() string {
	return e.Message
}

func ArgsParser(parameters []string) ([]string, error) {
	inputs := parameters
	if len(inputs) == 0 {
		inputs = os.Args
	}

	if len(inputs) < 6 || len(inputs) > 6 {
		log.Fatalln("Usage: ./main.go {filename}.h264 {bytesToRemove} {offset} {ratio} {reverse}\n" +
			"Parameters:\n" +
			"  filename      string  : Input file name with .h264 extension\n" +
			"  bytesToRemove int     : 0 to 100 if ratio is true, 0 to n if ratio is false\n" +
			"  offset        int     : 1 to 100 Offset position for removal\n" +
			"  ratio         bool 	 : Ratio for processing (true: 1/false: 0)\n" +
			"  reverse       bool    : Reverse the operation (true: 1/false: 0)")
	}

	// validate the input parameters
	if isInvalidFileName(inputs[1]) ||
		isInvalidBytesToRemoveandRatio(inputs[2], inputs[4]) ||
		isInvalidOffset(inputs[3]) ||
		isInvalidReverse(inputs[5]) {

		return inputs, &argsError{Message: "Invalid Prameter. Please check the input format and try again."}
	}

	return inputs, nil
}

func isInvalidFileName(filename string) bool {
	if strings.HasSuffix(filename, ".h264") {
		return false
	}
	fmt.Println("Please check the file foramt value and try again.")
	return true
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
	} else if ratioInt == 0 {
		if bytesToRemoveInt < 0 {
			fmt.Println("Please check the bytesToRemove value and try again.\n" +
				"ratio is false, but bytesToRemove is not in the range of 0 to n.")
			return true
		}
	} else {
		fmt.Println("Please check the ratio value and try again.\n" +
			"ratio is only 1(true) or 0(false), but is not.")
		return true
	}
	return false
}

func isInvalidOffset(offset string) bool {
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		fmt.Println("Please check the offset value and try again.\nError: ", err)
		return true
	}
	if offsetInt < 1 || offsetInt > 100 {
		fmt.Println("Offset must be in the range of 1 to 100.")
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
