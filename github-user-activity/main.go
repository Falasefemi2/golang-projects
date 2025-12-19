package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	DefaultBaseURL    = "https://api.github.com"
	DefaultEventLimit = 10
	DefaultTimeout    = 30 * time.Second
	GithubRateLimit   = 60

	CmdEvents = "events"
	CmdStats  = "stats"
	CmdRepos  = "repos"
	CmdHelp   = "help"

	SortUpdated = "updated"
	SortStars   = "stars"
	SortForks   = "forks"
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
	BaseUrl      string
	httpClient   *http.Client
	Subcommands []Subcommand
}

type Subcommand struct {
	Name      string
	Run       func(args []string) error
	Usage     string
	RequiresAuth bool
}

func (gc *GithubClient) Run(args []string) error {
	if len(args) < 1 {
		printUsage()
		return nil
	}

	cmdName := args[0]
	for _, cmd := range gc.Subcommands {
		if cmd.Name == cmdName {
			return cmd.Run(args[1:])
		}
	}

	fmt.Println("unknown command:", cmdName)
	printUsage()
	return nil
}

func main() {
	client := NewGithubClient(DefaultBaseURL, DefaultTimeout)
	if err := client.Run(os.Args[1:]); err != nil {
		handleError(err)
		os.Exit(1)
	}
}

func NewGithubClient(baseURL string, timeout time.Duration) *GithubClient {
	client := &GithubClient{
		BaseUrl: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
	client.Subcommands = []Subcommand{
		{
			Name:      CmdEvents,
			Run:       client.eventsCmd,
			Usage:     "events <username> [limit]",
			RequiresAuth: false,
		},
		{
			Name:      CmdStats,
			Run:       client.statsCmd,
			Usage:     "stats <username>",
			RequiresAuth: false,
		},
		{
			Name:      CmdRepos,
			Run:       client.reposCmd,
			Usage:     "repos <username> [sort] [limit]",
			RequiresAuth: false,
		},
		{
			Name:      CmdHelp,
			Run:       helpCmd,
			Usage:     "help",
			RequiresAuth: false,
		},
	}
	return client
}

func (gc *GithubClient) eventsCmd(args []string) error {
	if len(args) < 1 {
		return errors.New("username is required")
	}
	username := args[0]
	limit := DefaultEventLimit

	if len(args) > 1 {
		_, err := fmt.Sscanf(args[1], "%d", &limit)
		if err != nil || limit < 1 || limit > 300 {
			return errors.New("limit must be a number between 1 and 300")
		}
	}

	events, err := gc.GetEvents(username)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		fmt.Println("No events found for user:", username)
		return nil
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
	return nil
}

func (gc *GithubClient) statsCmd(args []string) error {
	if len(args) < 1 {
		return errors.New("username is required")
	}
	username := args[0]

	var user User
	var repos []Repo
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		u, err := gc.GetUserInfo(username)
		if err != nil {
			errChan <- err
			return
		}
		user = *u
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		r, err := gc.GetRepoStat(username, SortUpdated, 10)
		if err != nil {
			errChan <- err
			return
		}
		repos = r
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
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

	fmt.Printf("\n📚 Top 10 Repositories for @%s (sorted by updated)\n", username)
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

	return nil
}

func (gc *GithubClient) reposCmd(args []string) error {
	if len(args) < 1 {
		return errors.New("username is required")
	}
	username := args[0]
	sort := "updated"
	limit := 10

	if len(args) > 1 {
		sort = args[1]
		if sort != SortUpdated && sort != SortStars && sort != SortForks {
			return errors.New("sort must be 'updated', 'stars', or 'forks'")
		}
	}
	if len(args) > 2 {
		_, err := fmt.Sscanf(args[2], "%d", &limit)
		if err != nil || limit < 1 || limit > 100 {
			return errors.New("limit must be a number between 1 and 100")
		}
	}
	repos, err := gc.GetRepoStat(username, sort, limit)
	if err != nil {
		return err
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
	return nil
}

func helpCmd(args []string) error {
	printUsage()
	return nil
}

func printUsage() {
	fmt.Println(`
GitHub CLI Wrapper
A simple command-line tool for querying GitHub user activity.

Usage:
  github-cli <command> <username> [options]

Commands:
  events <username> [limit]    Show recent public events for a user
  stats <username>             Show summary of user stats
  repos <username> [sort] [limit] Show user repositories
  help                         Show this help message

Examples:
  github-cli events torvalds
  github-cli events torvalds 5
  github-cli stats torvalds
  github-cli repos torvalds stars
  github-cli help

Notes:
  - <username> is a GitHub username
  - [limit] is optional (default: 10, max: 300)
  - Unauthenticated requests are limited to 60 per hour
`)
}

func (gc *GithubClient) doRequest(url string, result interface{}) error {
	resp, err := gc.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to GitHub API: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	// Continue
	case http.StatusNotFound:
		return fmt.Errorf("user not found on GitHub")
	case http.StatusForbidden:
		return errors.New("rate limit exceeded. GitHub allows 60 unauthenticated requests per hour. Use a personal access token for higher limits")
	case http.StatusUnauthorized:
		return errors.New("authentication failed. Please check your credentials")
	case http.StatusInternalServerError, http.StatusServiceUnavailable:
		return errors.New("GitHub API is temporarily unavailable. Please try again later")
	default:
		return fmt.Errorf("GitHub API error (HTTP %d): %s", resp.StatusCode, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to parse API response: %w", err)
	}
	return nil
}

func (gc *GithubClient) GetEvents(username string) ([]Event, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	url := fmt.Sprintf("%s/users/%s/events", gc.BaseUrl, username)
	var events []Event
	err := gc.doRequest(url, &events)
	return events, err
}

func (gc *GithubClient) GetUserInfo(username string) (*User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	url := fmt.Sprintf("%s/users/%s", gc.BaseUrl, username)
	var user User
	err := gc.doRequest(url, &user)
	return &user, err
}

func (gc *GithubClient) GetRepoStat(username string, sort string, pageNum int) ([]Repo, error) {
	if username == "" || sort == "" {
		return nil, errors.New("username and sort cannot be empty")
	}
	url := fmt.Sprintf("%s/users/%s/repos?sort=%s&per_page=%d", gc.BaseUrl, username, sort, pageNum)
	var repos []Repo
	err := gc.doRequest(url, &repos)
	return repos, err
}

func handleError(err error) {
	fmt.Printf("Error: %v\n", err)
	if errors.Is(err, errors.New("rate limit exceeded")) {
		fmt.Println("\nTip: Generate a personal access token at https://github.com/settings/tokens")
		fmt.Println("Then set it as an environment variable: export GITHUB_TOKEN=your_token")
	}
}

