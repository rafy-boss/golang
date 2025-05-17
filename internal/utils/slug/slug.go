package slug

import (
	"regexp"
	"strings"
)

func GenerateSlug(input string) string {

	return strings.ToLower(strings.Trim(regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(input, "-"), "-"))
}