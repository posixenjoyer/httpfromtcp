package request

import (
	"slices"
	"strconv"
	"strings"
	"unicode"
)

func verifyMethod(methodName string) bool {
	return verifyMethodName(methodName) && verifyCaseAndContents(methodName)
}

func verifyMethodName(methodName string) bool {
	return slices.Contains(httpMethods, methodName)
}

func verifyCaseAndContents(methodName string) bool {
	for _, r := range methodName {
		if unicode.IsUpper(r) && unicode.IsLetter(r) {
			continue
		}
		return false
	}

	return true
}

func verifyVersion(version string) bool {
	if !strings.Contains(version, "/") {
		return false
	}
	parts := strings.Split(version, "/")
	if len(parts) != 2 {
		return false
	}
	protocol := parts[0]
	versionNum := parts[1]
	if protocol != "HTTP" {
		return false
	}
	versionNumParts := strings.Split(versionNum, ".")
	if len(versionNumParts) != 2 {
		return false
	}
	major := versionNumParts[0]
	minor := versionNumParts[1]
	if _, err := strconv.Atoi(major); err != nil {
		return false
	}
	if _, err := strconv.Atoi(minor); err != nil {
		return false
	}
	return true
}
