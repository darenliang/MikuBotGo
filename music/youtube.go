package music

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/guregu/null.v4"
	"os/exec"
	"strings"
)

type (
	// Query type enum
	QUERYTYPE int

	// Empty youtube receiver
	Youtube struct{}

	SongResponse struct {
		URL        string      `json:"url"`
		UploadDate null.String `json:"upload_date"`
		Duration   null.Float  `json:"duration"`
		Thumbnail  string      `json:"thumbnail"`
		Title      string      `json:"title"`
	}

	PlaylistResponse struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	}
)

const (
	ERRORTYPE QUERYTYPE = iota
	SONGTYPE
	PLAYLISTTYPE
)

func (youtube Youtube) getType(input string) QUERYTYPE {
	firstLine := input[:strings.Index(input, "\n")]
	var firstLineMap map[string]interface{}
	if err := json.Unmarshal([]byte(firstLine), &firstLineMap); err != nil {
		return ERRORTYPE
	}

	// Inspect json result
	if _, ok := firstLineMap["_type"]; ok {
		return PLAYLISTTYPE
	}
	if _, ok := firstLineMap["_filename"]; ok {
		return SONGTYPE
	}

	return ERRORTYPE
}

func (youtube Youtube) YoutubeDLLink(input string) (QUERYTYPE, *string, error) {
	cmd := exec.Command("youtube-dl", "--default-search", "ytsearch", "--skip-download", "--print-json", "--flat-playlist", "--format", "best", "-x", input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ERRORTYPE, nil, err
	}
	str := out.String()
	return youtube.getType(str), &str, nil
}

func (youtube Youtube) YoutubeDLQuery(input string) (QUERYTYPE, *string, error) {
	cmd := exec.Command("youtube-dl", fmt.Sprintf("ytsearch4:\"%s\"", input), "--skip-download", "--print-json", "--flat-playlist")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ERRORTYPE, nil, err
	}
	str := out.String()
	return PLAYLISTTYPE, &str, nil
}

func (youtube Youtube) GetSong(input string) (*SongResponse, error) {
	var resp SongResponse
	err := json.Unmarshal([]byte(input), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (youtube Youtube) GetPlaylist(input string) (*[]PlaylistResponse, error) {
	lines := strings.Split(input, "\n")
	videos := make([]PlaylistResponse, 0)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var video PlaylistResponse
		err := json.Unmarshal([]byte(line), &video)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return &videos, nil
}
