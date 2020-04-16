package framework

import (
	"encoding/json"
	"github.com/paul-mannino/go-fuzzywuzzy"
	"strings"
)

type Openings []struct {
	Title  string `json:"title"`
	Source string `json:"source"`
	File   string `json:"file"`
}

// Filter returns a new slice containing all openings in the slice that satisfy the predicate f
func Filter(vs Openings, f func(string) bool) Openings {
	vsf := make(Openings, 0)
	for _, v := range vs {
		if f(v.Title) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// isOpening returns whether or not string is an opening
func IsOpening(str string) bool {
	return strings.Contains(str, "Opening")
}

// Return openings
func GetOpenings() Openings {
	result := GetJsonString("https://openings.moe/api/list.php")
	tmp := Openings{}
	_ = json.Unmarshal([]byte(result), &tmp)
	return Filter(tmp, IsOpening)
}

// GetStringValidation checks if string is correct
func GetStringValidation(answers []string, guess string) bool {
	for _, val := range answers {
		if fuzzy.TokenSortRatio(val, guess) > 60 {
			return true
		}
	}
	return false
}
