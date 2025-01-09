package tools

import (
	"strconv"
	"strings"
)

func ExtractID(path string) (uint, error) {

	parts := strings.Split(path, "/")
	id, err := strconv.Atoi(parts[len(parts)-1])

	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
