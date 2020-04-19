package framework

type MALAnime struct {
	URL           string   `json:"url"`
	ImageURL      string   `json:"image_url"`
	OpeningThemes []string `json:"opening_themes"`
	EndingThemes  []string `json:"ending_themes"`
}
