package framework

import (
	"encoding/json"
	"github.com/paul-mannino/go-fuzzywuzzy"
	"io/ioutil"
	"math/rand"
	"time"
)

type Openings []struct {
	Name  string `json:"name"`
	Songs []struct {
		Songname string `json:"songname"`
		URL      string `json:"url"`
	} `json:"songs"`
}

func init() {
	// Generate random seed
	rand.Seed(time.Now().UnixNano())
}

// Return openings
func GetOpenings() Openings {
	file, _ := ioutil.ReadFile("data/dataset_filtered.json")
	tmp := Openings{}
	_ = json.Unmarshal(file, &tmp)
	return tmp
}

// GetStringValidation checks if string is correct
func GetStringValidation(answers []string, guess string) bool {
	for _, val := range answers {
		if fuzzy.TokenSortRatio(val, guess) > 70 {
			return true
		}
	}
	return false
}
