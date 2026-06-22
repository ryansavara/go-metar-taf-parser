package metartafparser

import (
	"regexp"
	"strings"
)

var fractionalVisRe = regexp.MustCompile(`^\d+/\d+SM$`)
var simpleNumRe = regexp.MustCompile(`^\d+$`)

func tokenize(input string) []string {
	parts := splitOnEquals(input)
	var result []string
	for _, part := range parts {
		tokens := splitAndMergeVisibility(part)
		result = append(result, tokens...)
	}
	return result
}

func splitOnEquals(input string) []string {
	parts := strings.Split(input, "=")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func splitAndMergeVisibility(input string) []string {
	raw := strings.Fields(input)
	if len(raw) == 0 {
		return nil
	}
	var merged []string
	for i := range raw {
		if i > 0 && fractionalVisRe.MatchString(raw[i]) && simpleNumRe.MatchString(raw[i-1]) {
			merged[len(merged)-1] = merged[len(merged)-1] + " " + raw[i]
		} else {
			merged = append(merged, raw[i])
		}
	}
	return merged
}
