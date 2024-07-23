package apputil

import (
	"regexp"
	"strings"
)

// RewriteRulesRegex applies re-write rules
// eg. "/app/*/*": "/$2"
func RewriteRulesRegex(rewrite map[string]string) map[*regexp.Regexp]string {
	// Initialize
	rulesRegex := map[*regexp.Regexp]string{}
	for k, v := range rewrite {
		k = regexp.QuoteMeta(k)
		k = strings.ReplaceAll(k, `\*`, "(.*?)")
		if strings.HasPrefix(k, `\^`) {
			k = strings.ReplaceAll(k, `\^`, "^")
		}
		k += "$"
		rulesRegex[regexp.MustCompile(k)] = v
	}
	return rulesRegex
}
