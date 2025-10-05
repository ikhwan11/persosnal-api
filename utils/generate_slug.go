package utils

import (
	"regexp"
	"strings"
)

func GenerateSlug(name string) string {
	slug := strings.ToLower(name)

	re := regexp.MustCompile(`[^a-z0-9!?\. ]+`)
	slug = re.ReplaceAllString(slug, "")

	slug = strings.ReplaceAll(slug, " ", "-")

	re2 := regexp.MustCompile(`-+`)
	slug = re2.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}
