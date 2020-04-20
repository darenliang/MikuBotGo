package framework

import (
	"github.com/paul-mannino/go-fuzzywuzzy"
	"strings"
)

// GetStringValidation checks if string is correct
func GetStringValidation(answers []string, guess string) bool {
	guess = strings.ToLower(guess)
	for _, val := range answers {
		val = strings.ToLower(val)
		if fuzzy.TokenSortRatio(val, guess) > 65 {
			return true
		}
	}
	return false
}
