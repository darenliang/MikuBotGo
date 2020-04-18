package framework

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

type Openings []struct {
	Title  string `json:"title"`
	Source string `json:"source"`
	File   string `json:"file"`
}

type Openings2 []struct {
	Year   int `json:"year"`
	Animes []struct {
		Name  string   `json:"name"`
		Songs []string `json:"songs"`
	} `json:"animes"`
}

type Openings3 []struct {
	Name  string   `json:"name"`
	Songs []string `json:"songs"`
}

func init() {
	// Generate random seed
	rand.Seed(time.Now().UnixNano())
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
	return tmp
	// return Filter(tmp, IsOpening)
}

// Return openings version 2
func GetOpenings2() Openings2 {
	file, _ := ioutil.ReadFile("data/dataset.json")
	tmp := Openings2{}
	_ = json.Unmarshal(file, &tmp)
	return tmp
}

// Return openings version 3
func GetOpenings3() Openings3 {
	file, _ := ioutil.ReadFile("data/dataset_filtered.json")
	tmp := Openings3{}
	_ = json.Unmarshal(file, &tmp)
	return tmp
}

func GetRandomYear() int {
	rank := rand.Int() % 231
	yearOffset := 21
	rank -= yearOffset
	for rank >= 0 {
		yearOffset--
		rank -= yearOffset
	}
	return yearOffset + 1999
}

// GetStringValidation checks if string is correct
func GetStringValidation(answers []string, guess string) bool {
	for _, val := range answers {
		if val == guess {
			return true
		}
	}
	return false
}
