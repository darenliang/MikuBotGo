package framework

import (
	"fmt"
	"net/url"
)

type AniListAnimeSearchResponse struct {
	Data struct {
		Media struct {
			IDMal int `json:"idMal"`
			Title struct {
				Romaji        string `json:"romaji"`
				English       string `json:"english"`
				Native        string `json:"native"`
				UserPreferred string `json:"userPreferred"`
			} `json:"title"`
			Type        string `json:"type"`
			Format      string `json:"format"`
			Status      string `json:"status"`
			Description string `json:"description"`
			StartDate   struct {
				Year  int `json:"year"`
				Month int `json:"month"`
				Day   int `json:"day"`
			} `json:"startDate"`
			EndDate struct {
				Year  int `json:"year"`
				Month int `json:"month"`
				Day   int `json:"day"`
			} `json:"endDate"`
			Season     string `json:"season"`
			SeasonYear int    `json:"seasonYear"`
			Episodes   int    `json:"episodes"`
			Source     string `json:"source"`
			CoverImage struct {
				ExtraLarge string `json:"extraLarge"`
				Color      string `json:"color"`
			} `json:"coverImage"`
			Genres       []string `json:"genres"`
			AverageScore int      `json:"averageScore"`
			Studios      struct {
				Edges []struct {
					Node struct {
						IsAnimationStudio bool   `json:"isAnimationStudio"`
						Name              string `json:"name"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"studios"`
			Rankings []struct {
				Rank    int    `json:"rank"`
				Type    string `json:"type"`
				AllTime bool   `json:"allTime"`
			} `json:"rankings"`
			SiteURL string `json:"siteUrl"`
		} `json:"Media"`
	} `json:"data"`
}

type AniListMangaSearchResponse struct {
	Data struct {
		Media struct {
			IDMal int `json:"idMal"`
			Title struct {
				Romaji        string `json:"romaji"`
				English       string `json:"english"`
				Native        string `json:"native"`
				UserPreferred string `json:"userPreferred"`
			} `json:"title"`
			Type        string `json:"type"`
			Format      string `json:"format"`
			Status      string `json:"status"`
			Description string `json:"description"`
			StartDate   struct {
				Year  int `json:"year"`
				Month int `json:"month"`
				Day   int `json:"day"`
			} `json:"startDate"`
			EndDate struct {
				Year  int `json:"year"`
				Month int `json:"month"`
				Day   int `json:"day"`
			} `json:"endDate"`
			Volumes    int    `json:"volumes"`
			Chapters   int    `json:"chapters"`
			Source     string `json:"source"`
			CoverImage struct {
				ExtraLarge string `json:"extraLarge"`
				Color      string `json:"color"`
			} `json:"coverImage"`
			Staff struct {
				Edges []struct {
					Node struct {
						Name struct {
							Full string `json:"full"`
						} `json:"name"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"staff"`
			Genres       []string `json:"genres"`
			AverageScore int      `json:"averageScore"`
			Rankings     []struct {
				Rank    int    `json:"rank"`
				Type    string `json:"type"`
				AllTime bool   `json:"allTime"`
			} `json:"rankings"`
			SiteURL string `json:"siteUrl"`
		} `json:"Media"`
	} `json:"data"`
}

type TraceSearchResult struct {
	Data struct {
		Media struct {
			Title struct {
				UserPreferred string `json:"userPreferred"`
			} `json:"title"`
			CoverImage struct {
				ExtraLarge string `json:"extraLarge"`
				Color      string `json:"color"`
			} `json:"coverImage"`
		} `json:"Media"`
	} `json:"data"`
}

func AnilistAnimeSearchQuery(query string) map[string]string {
	jsonQuery := map[string]string{
		"query": fmt.Sprintf(`
{
  Media(search: "%s", type: ANIME) {
    idMal
    title {
      romaji
      english
      native
      userPreferred
    }
    type
    format
    status
    description
    startDate {
      year
      month
      day
    }
    endDate {
      year
      month
      day
    }
    season
    seasonYear
    episodes
    source
    coverImage {
      extraLarge
      color
    }
    genres
    averageScore
    studios {
      edges {
        node {
          isAnimationStudio
          name
        }
      }
    }
    rankings {
      rank
      type
      allTime
    }
    siteUrl
  }
}
`, url.QueryEscape(query))}
	return jsonQuery
}

func AnilistAnimeIDQuery(id int) map[string]string {
	jsonQuery := map[string]string{
		"query": fmt.Sprintf(`
{
  Media(id: %d) {
    title {
      userPreferred
    }
    coverImage {
      extraLarge
      color
    }
  }
}

`, id)}
	return jsonQuery
}

func AnilistMangaSearchQuery(query string) map[string]string {
	jsonQuery := map[string]string{
		"query": fmt.Sprintf(`
{
  Media(search: "%s", type: MANGA) {
    idMal
    title {
      romaji
      english
      native
      userPreferred
    }
    type
    format
    status
    description
    startDate {
      year
      month
      day
    }
    endDate {
      year
      month
      day
    }
    volumes
    chapters
    source
    coverImage {
      extraLarge
      color
    }
    staff {
      edges {
        node {
          name {
            full
          }
        }
      }
    }
    genres
    averageScore
    rankings {
      rank
      type
      allTime
    }
    siteUrl
  }
}`, url.QueryEscape(query))}
	return jsonQuery
}
