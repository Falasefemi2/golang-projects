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

type User struct {
	Login           string    `json:"login"`
	ID              int       `json:"id"`
	URL             string    `json:"url"`
	Hireable        bool      `json:"hireable"`
	Bio             string    `json:"bio"`
	TwitterUsername string    `json:"twitter_username"`
	PublicRepos     int       `json:"public_repos"`
	PublicGist      int       `json:"public_gist"`
	Followers       int       `json:"followers"`
	Following       int       `json:"following"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Repo struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Private         bool      `json:"private"`
	Description     string    `json:"description"`
	URL             string    `json:"url"`
	HtmlURL         string    `json:"html_url"`
	Language        string    `json:"language"`
	StargazersCount int       `json:"stargazers_count"`
	WatchersCount   int       `json:"watchers_count"`
	ForksCount      int       `json:"forks_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
	Topics          []string  `json:"topics"`
	Owner           struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		HtmlURL           string `json:"html_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
	} `json:"owner"`
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

	case "stats":
		if len(os.Args) < 3 {
			fmt.Println("Error: username is required")
			printUsage()
			return
		}
		username := os.Args[2]
		client := NewGithubClient(DefaultBaseURL, DefaultTimeout)
		user, err := client.GetUserInfo(username)
		if err != nil {
			handleError(err)
			return
		}

		fmt.Printf("\n📊 GitHub Stats for @%s\n", user.Login)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		if user.Bio != "" {
			fmt.Printf("Bio: %s\n", user.Bio)
		}

		yearsOnGithub := time.Since(user.CreatedAt).Hours() / 24 / 365
		fmt.Printf("Member since: %s (%.1f years)\n", user.CreatedAt.Format("2006-01-02"), yearsOnGithub)

		fmt.Printf("Public Repositories: %d\n", user.PublicRepos)
		fmt.Printf("Public Gists: %d\n", user.PublicGist)
		fmt.Printf("Followers: %d\n", user.Followers)
		fmt.Printf("Following: %d\n", user.Following)

		if user.TwitterUsername != "" {
			fmt.Printf("Twitter: @%s\n", user.TwitterUsername)
		}

		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	case "repos":
		if len(os.Args) < 3 {
			fmt.Println("Error: username is required")
			fmt.Println("Usage: github-cli repos <username> [sort] [limit]")
			fmt.Println("Sort options: updated (default), stars, forks")
			printUsage()
			return
		}
		username := os.Args[2]
		sort := "updated"
		limit := 10

		if len(os.Args) > 3 {
			sort = os.Args[3]
			if sort != "updated" && sort != "stars" && sort != "forks" {
				fmt.Println("Error: sort must be 'updated', 'stars', or 'forks'")
				return
			}
		}
		if len(os.Args) > 4 {
			_, err := fmt.Sscanf(os.Args[4], "%d", &limit)
			if err != nil || limit < 1 || limit > 100 {
				fmt.Println("Error: limit must be a number between 1 and 100")
				return
			}
		}
		client := NewGithubClient(DefaultBaseURL, DefaultTimeout)
		repos, err := client.GetRepoStat(username, sort, limit)
		if err != nil {
			handleError(err)
			return
		}
		fmt.Printf("\n📚 Top %d Repositories for @%s (sorted by %s)\n", len(repos), username, sort)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		for i, repo := range repos {
			fmt.Printf("\n%d. %s\n", i+1, repo.Name)

			if repo.Description != "" {
				fmt.Printf("   Description: %s\n", repo.Description)
			}

			if repo.Language != "" {
				fmt.Printf("   Language: %s\n", repo.Language)
			}

			fmt.Printf("   ⭐ Stars: %d  🍴 Forks: %d\n", repo.StargazersCount, repo.ForksCount)
			fmt.Printf("   Last updated: %s\n", repo.UpdatedAt.Format("2006-01-02"))
			fmt.Printf("   Link: %s\n", repo.HtmlURL)
		}

		fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
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

func (gc *GithubClient) GetUserInfo(username string) (*User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	url := fmt.Sprintf("%s/users/%s", gc.BaseUrl, username)
	client := &http.Client{
		Timeout: gc.Timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Github SPI: %w", err)
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
	var users User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}
	return &users, nil
}

func (gc *GithubClient) GetRepoStat(username string, sort string, pageNum int) ([]Repo, error) {
	if username == "" || sort == "" {
		return nil, errors.New("username cannot be empty | sort cannot be empty")
	}
	url := fmt.Sprintf("%s/users/%s/repos?sort=%s&per_page=%d", gc.BaseUrl, username, sort, pageNum)
	client := &http.Client{
		Timeout: gc.Timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Github API: %w", err)
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
	var repos []Repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}
	return repos, nil
}

func handleError(err error) {
	fmt.Printf("Error: %v\n", err)
	if errors.Is(err, errors.New("rate limit exceeded")) {
		fmt.Println("\nTip: Generate a personal access token at https://github.com/settings/tokens")
		fmt.Println("Then set it as an environment variable: export GITHUB_TOKEN=your_token")
	}
}
