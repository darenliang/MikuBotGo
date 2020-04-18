package framework

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// Index returns index of string found in list. Return -1 if not found
func Index(str string, vs []string) int {
	for i, val := range vs {
		if val == str {
			return i
		}
	}
	return -1
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	err = resp.Body.Close()

	if err != nil {
		return err
	}

	err = out.Close()

	if err != nil {
		return err
	}

	return err
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

// GeneratePreviewDesc returns preview description
// Between space after : and \n
func GeneratePreviewDesc(value string) string {
	posFirst := strings.Index(value, ":")
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, "\n")
	if posLast == -1 {
		return ""
	}
	return value[posFirst+2 : posLast]
}

func ParseDate(year, month, day int) string {
	t, _ := time.Parse("2006-01-02",
		fmt.Sprintf("%d-%02d-%02d", year, month, day))
	return t.UTC().Format("January 2, 2006")
}
