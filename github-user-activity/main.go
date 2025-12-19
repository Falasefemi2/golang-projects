package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	DefaultBaseURL    = "https://api.github.com"
	DefaultEventLimit = 10
	DefaultTimeout    = 30 * time.Second
	GithubRateLimit   = 60
)

type Event struct {
	ID   string `json:"id"`
	Type string `json:"type"`

	Actor struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
	} `json:"actor"`

	Repo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"repo"`

	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
}

type GithubClient struct {
	BaseUrl string
	Timeout time.Duration
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "events":
		if len(os.Args) < 3 {
			fmt.Println("Error: username is required")
			printUsage()
			return
		}
		username := os.Args[2]
		limit := DefaultEventLimit

		if len(os.Args) > 3 {
			_, err := fmt.Sscanf(os.Args[3], "%d", &limit)
			if err != nil || limit < 1 || limit > 300 {
				fmt.Println("Error: limit must be a number between 1 and 300")
				return
			}
		}

		client := NewGithubClient(DefaultBaseURL, DefaultTimeout)
		events, err := client.GetEvents(username)
		if err != nil {
			handleError(err)
			return
		}
		if len(events) == 0 {
			fmt.Println("No events found for user:", username)
			return
		}

		displayCount := limit
		if displayCount > len(events) {
			displayCount = len(events)
		}
		fmt.Printf("Showing %d of %d events for user '%s':\n\n", displayCount, len(events), username)
		for i := 0; i < displayCount; i++ {
			event := events[i]
			fmt.Printf("%d. [%s] %s\n", i+1, event.CreatedAt.Format("2006-01-02 15:04"), event.Type)
			fmt.Printf("   Repository: %s\n", event.Repo.Name)
			fmt.Printf("   Actor: %s\n\n", event.Actor.Login)
		}
	case "help":
		printUsage()
	default:
		fmt.Println("unknown command:", command)
		printUsage()
	}
}

func NewGithubClient(baseURL string, timeout time.Duration) *GithubClient {
	return &GithubClient{
		BaseUrl: baseURL,
		Timeout: timeout,
	}
}

func printUsage() {
	fmt.Println(`
GitHub CLI Wrapper
A simple command-line tool for querying GitHub user activity.

Usage:
  github-cli <command> <username> [options]

Commands:
  events <username> [limit]    Show recent public events for a user
  help                         Show this help message

Examples:
  github-cli events torvalds
  github-cli events torvalds 5
  github-cli help

Notes:
  - <username> is a GitHub username
  - [limit] is optional (default: 10, max: 300)
  - Unauthenticated requests are limited to 60 per hour
`)
}

func (gc *GithubClient) GetEvents(username string) ([]Event, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	url := fmt.Sprintf("%s/users/%s/events", gc.BaseUrl, username)
	client := &http.Client{
		Timeout: gc.Timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to GitHub API: %w", err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("user '%s' not found on GitHub", username)
	case http.StatusForbidden:
		return nil, errors.New("rate limit exceeded. GitHub allows 60 unauthenticated requests per hour. Use a personal access token for higher limits")
	case http.StatusUnauthorized:
		return nil, errors.New("authentication failed. Please check your credentials")
	case http.StatusInternalServerError, http.StatusServiceUnavailable:
		return nil, errors.New("GitHub API is temporarily unavailable. Please try again later")
	default:
		return nil, fmt.Errorf("GitHub API error (HTTP %d): %s", resp.StatusCode, resp.Status)
	}
	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}
	return events, nil
}

func handleError(err error) {
	fmt.Printf("Error: %v\n", err)

	// Provide helpful suggestions based on error type
	if errors.Is(err, errors.New("rate limit exceeded")) {
		fmt.Println("\nTip: Generate a personal access token at https://github.com/settings/tokens")
		fmt.Println("Then set it as an environment variable: export GITHUB_TOKEN=your_token")
	}
}
