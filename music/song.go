package music

import (
	"os/exec"
	"strconv"
)

type Song struct {
	Media    string
	Title    string
	Duration *string
	Id       string
}

func (song Song) Ffmpeg() *exec.Cmd {
	return exec.Command("ffmpeg", "-reconnect", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "10",
		"-i", song.Media, "-f", "s16le", "-ar", strconv.Itoa(FRAMERATE), "-ac",
		strconv.Itoa(CHANNELS), "pipe:1")
}

func NewSong(media, title, id string) *Song {
	song := new(Song)
	song.Media = media
	song.Title = title
	song.Id = id
	return song
}
