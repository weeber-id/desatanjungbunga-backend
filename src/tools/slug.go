package tools

import (
	"regexp"
	"strings"
)

// GenerateSlug from article title
func GenerateSlug(title string) (string, error) {

	// filter only get alphanum and space
	refilter, err := regexp.Compile(`[0-9A-z ]+`)
	if err != nil {
		return "", err
	}
	filter1 := strings.Join(refilter.FindAllString(title, -1), "")

	// replace space by using "-"
	respace, err := regexp.Compile(` `)
	if err != nil {
		return "", err
	}
	filter2 := respace.ReplaceAllString(filter1, "-")
	return filter2, nil
}
