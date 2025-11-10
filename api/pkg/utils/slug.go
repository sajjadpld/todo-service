package utils

import "strings"

func Slug(phrase string) (slug string) {
	slug = strings.ToLower(phrase)
	slug = strings.ReplaceAll(slug, " ", "-")
	return
}
