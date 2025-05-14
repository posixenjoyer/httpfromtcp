package request

import (
	"slices"
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
	return version == "HTTP/1.1"
}
