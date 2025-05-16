package headers

import (
	"fmt"
	"strings"
	"unicode"
)

func validateFieldName(field string) bool {
	fmt.Printf("field-name: %s\n", field)

	for n, c := range field {
		if unicode.IsDigit(c) || unicode.IsLetter(c) || strings.Contains(validFieldSpecial, string(c)) {
			continue
		} else {
			fmt.Printf("Returning false: %c @ Index: %d\n", c, n)
			return false
		}
	}

	return true
}
