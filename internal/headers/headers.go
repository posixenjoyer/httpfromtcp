package headers

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func (h Headers) Parse(headerLine []byte) (int, bool, error) {

	done := false

	crlfIndex := strings.Index(string(headerLine), "\r\n")
	if crlfIndex < 0 {
		return 0, false, nil
	}

	if crlfIndex == 0 {
		return 0, true, nil
	}

	count := crlfIndex + len("\r\n")

	cleanHeaderField := strings.TrimSpace(string(headerLine[:crlfIndex]))

	keyPair := strings.SplitN(cleanHeaderField, ":", 2)

	if len(keyPair) != 2 {
		return count, false, errors.New("no valid key found")
	}

	keyName := strings.ToLower(keyPair[0])
	valName := strings.TrimSpace(keyPair[1])

	if unicode.IsSpace(rune(keyName[len(keyName)-1])) {
		return 0, false, errors.New("invalid header name format")
	}

	if !validateFieldName(keyName) {
		fmt.Printf("validation failed")
		return count, false, errors.New("invalid character in field-name")
	}

	if ok := h[keyName]; ok != "" {
		h[keyName] = h[keyName] + ", " + valName
	} else {
		h[keyName] = valName
	}
	return count, done, nil
}
