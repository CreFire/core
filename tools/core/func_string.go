package core

import "strings"

func StrSplitAll(s string, sep string) []string {
	return strings.Split(s, sep)
}

func StrReplaceAll(s string, old string, new string) string {
	return strings.Replace(s, old, new, -1)
}

func StrSplitDefault(s string, delimiters []string) []string {
	if delimiters != nil {
		delimiters = []string{",", "`", " ", "、", ":"}
	}
	for _, delimiter := range delimiters {
		s = strings.ReplaceAll(s, delimiter, " ")
	}
	return strings.Split(s, " ")
}
