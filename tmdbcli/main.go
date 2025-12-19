package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	apiBaseURL  = "https://api.themoviedb.org/3/movie"
	bearerToken = "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIyYWVjNmIzOWY1NzhiNThkNjRkNDU5ZWJkOThkNTFmMyIsIm5iZiI6MTY0MDAyOTI5OC4wMTksInN1YiI6IjYxYzBkYzcyNGRhM2Q0MDA2M2JlNGJiYyIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.bVaFCYoGl0pzbz-UpUKPTb3P_DVSMXj7tjREpYvS1j0"
)

type Movie struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIds         []int   `json:"genre_ids"`
	Id               int     `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	Popularity       float64 `json:"popularity"`
	PosterPath       string  `json:"poster_path"`
	ReleaseDate      string  `json:"release_date"`
	Title            string  `json:"title"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
}

type MovieResponse struct {
	Dates   map[string]string `json:"dates"`
	Page    int               `json:"page"`
	Results []Movie           `json:"results"`
}

func fetchMovies(endPoint string) (*MovieResponse, error) {
	url := fmt.Sprintf("%s/%s?language=en-US&page=1", apiBaseURL, endPoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API error: status %d - %s", res.StatusCode, string(body))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var movieResponse MovieResponse
	err = json.Unmarshal(body, &movieResponse)
	if err != nil {
		return nil, err
	}

	return &movieResponse, nil
}

func displayMovies(movieType string, response *MovieResponse) {
	fmt.Printf("\n=== %s Movies ===\n\n", strings.ToUpper(movieType))
	if len(response.Results) == 0 {
		fmt.Printf("No moies found")
		return
	}

	for i, movie := range response.Results {
		fmt.Printf("%d. %s (%s)\n", i+1, movie.Title, movie.ReleaseDate)
		fmt.Printf("   Rating: %.1f/10 (%d votes)\n", movie.VoteAverage, movie.VoteCount)
		fmt.Printf("   Popularity: %.2f\n", movie.Popularity)

		if movie.Overview != "" {
			overview := movie.Overview
			if len(overview) > 150 {
				overview = overview[:150] + "..."
			}
			fmt.Printf("Overview: %s\n", overview)
		}
		fmt.Println()
	}
}

func getMovies(movieType string) error {
	var endpoint string

	switch strings.ToLower(movieType) {
	case "playing":
		endpoint = "now_playing"
	case "popular":
		endpoint = "popular"
	case "top":
		endpoint = "top_rated"
	case "upcoming":
		endpoint = "upcoming"

	default:
		return fmt.Errorf("invalid type: %s. Use: playing, popular, top, or upcoming", movieType)
	}

	response, err := fetchMovies(endpoint)
	if err != nil {
		return err
	}

	displayMovies(movieType, response)
	return nil
}

func main() {
	movieType := flag.String("type", "", "Movie type: playing, popular, top, or upcoming")
	flag.Parse()

	if *movieType == "" {
		fmt.Println("Usage: tmdb-app --type <type>")
		fmt.Println("\nAvailable types:")
		fmt.Println("  playing   - Now playing movies")
		fmt.Println("  popular   - Popular movies")
		fmt.Println("  top       - Top rated movies")
		fmt.Println("  upcoming  - Upcoming movies")
		fmt.Println("\nExample: tmdb-app --type playing")
		os.Exit(1)
	}

	err := getMovies(*movieType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
