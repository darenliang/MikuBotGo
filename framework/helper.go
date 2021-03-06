package framework

import (
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"math/rand"
	"net/http"
	"time"
)

var HttpClient *http.Client

func init() {
	HttpClient = &http.Client{Timeout: 60 * time.Second}
}

// RandomString returns random string at fixed length
func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	buf := make([]rune, n)
	for i := range buf {
		buf[i] = letter[rand.Intn(len(letter))]
	}
	return string(buf)
}

func ParseDate(year, month, day int) string {
	t, _ := time.Parse("2006-01-02",
		fmt.Sprintf("%d-%02d-%02d", year, month, day))
	return t.UTC().Format("January 2, 2006")
}

// LoadImage from url
func LoadImage(url string) (image.Image, error) {
	// Get the data
	resp, err := HttpClient.Get(url)

	if err != nil {
		return nil, err
	}

	webImage, err := imaging.Decode(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	return webImage, nil
}

// UrlToStruct loads url json to struct target
func UrlToStruct(url string, target interface{}) error {
	resp, err := HttpClient.Get(url)

	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(target)

	errClose := resp.Body.Close()

	if errClose != nil {
		return errClose
	}

	if err != nil {
		return err
	}

	return nil
}

// Index returns index of str in arr
func Index(arr []string, str string) int {
	for i, v := range arr {
		if v == str {
			return i
		}
	}
	return -1
}
