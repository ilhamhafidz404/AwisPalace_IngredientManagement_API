package utils

import (
	"regexp"
	"strings"
)

func GenerateSlug(text string) string {
	slug := strings.ToLower(text)

	slug = strings.ReplaceAll(slug, " ", "-")

	reg, _ := regexp.Compile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	reg, _ = regexp.Compile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}
