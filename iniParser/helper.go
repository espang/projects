package main

import (
	"bytes"
	"fmt"
)

// ParseKeyValuePair expects a line number and a line like
// key=value
// and then returns key, value, nil
func ParseKeyValuePair(nbr int, line []byte) (string, string, error) {

	bytesSlices := bytes.Split(line, []byte{'='})

	if len(bytesSlices) != 2 {
		return "", "", fmt.Errorf("error in line %d. expect a sinlge '=' in the line: '%s'",
			nbr,
			line,
		)
	}

	key := string(bytes.TrimSpace(bytes.ToLower(bytesSlices[0])))
	value := string(bytes.TrimSpace(bytes.ToLower(bytesSlices[1])))

	return key, value, nil
}

// GetSectionFromLine returns a section name and true if the line
// is a section. otherwise will return "" and false
func GetSectionFromLine(line []byte) (string, bool) {
	trimmedLine := bytes.TrimSpace(line)
	if len(trimmedLine) < 2 {
		return "", false
	}
	if trimmedLine[0] == '[' && trimmedLine[len(trimmedLine)-1] == ']' {
		name := string(bytes.ToLower(trimmedLine[1 : len(trimmedLine)-1]))
		return name, name != ""
	}
	return "", false
}

func IsEmptyOrComment(line []byte) bool {
	trimmedLine := bytes.TrimSpace(line)
	return len(trimmedLine) == 0 || trimmedLine[0] == ';'
}
