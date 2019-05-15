package command

import "strings"

func getByteIndex(s string, b byte) int {
	index := strings.IndexByte(s, b)

	if index == -1 {
		return len(s)
	}

	return index
}
