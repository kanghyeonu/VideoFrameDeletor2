package util

import (
	"testing"
)

func TestArgsParser(t *testing.T) {
	tests := []struct {
		filename, bytesToRemove, offset, ratio, reverse string
		expectError                                     bool
	}{
		// 정상 케이스
		{"video.h264", "10", "1", "1", "0", false},
		{"foo.h264", "100", "50", "0", "1", false},
		{"poo.h264", "0", "100", "1", "1", false},

		// filename 검증
		{"video.mp4", "10", "1", "1", "0", true},
		{"video.h264", "10", "1", "1", "0", false},

		// bytesToRemove 검증
		{"video.h264", "150", "1", "1", "0", true},
		{"video.h264", "10", "1", "0", "0", false},

		// offset 검증
		{"video.h264", "10", "101", "1", "0", true},
		{"video.h264", "10", "0", "1", "0", true},
		{"video.h264", "10", "-1", "1", "0", true},

		// ratio 검증
		{"video.h264", "10", "15", "2", "0", true},
		{"video.h264", "10", "23", "0", "0", false},

		// byteToRemove & ratio 검증
		{"video.h264", "101", "15", "1", "0", true},
		{"video.h264", "-1", "15", "1", "0", true},
		{"video.h264", "-1", "15", "0", "0", true},

		// reverse 검증
		{"video.h264", "10", "11", "1", "2", true},
		{"video.h264", "10", "34", "1", "-1", true},
		{"video.h264", "10", "22", "1", "0", false},
	}

	for _, test := range tests {
		inputs := []string{"", test.filename, test.bytesToRemove, test.offset, test.ratio, test.reverse}
		_, err := ArgsParser(inputs)
		if (err != nil) != test.expectError {
			t.Errorf("ArgsParser(%v)\nerror = %v\nexpectError %v", inputs, err, test.expectError)
		}
	}
}
