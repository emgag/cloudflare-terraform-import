package util

import "strings"

func ToResourceName(i string) string {
	r := strings.NewReplacer(".", "-", "*", "wildcard")
	return strings.ToLower(r.Replace(i))
}
