package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
    "time"
)

// Repository represents a GitHub repository
type Repository struct {
    Name  string `json:"name"`
    Stars int    `json:"stargazers_count"`
}

// Fetch the list of repositories for a given GitHub user
func fetchRepositories(username string) ([]Repository, error) {
    var repos []Repository
    var nextPageURL string = "https://api.github.com/users/" + username + "/repos?per_page=100"

    for nextPageURL != "" {
        // Fetch the next page of repositories
        resp, err := http.Get(nextPageURL)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()

        // Read the response body
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }

        // Unmarshal the JSON response into a slice of Repository
        var pageRepos []Repository
        if err := json.Unmarshal(body, &pageRepos); err != nil {
            return nil, err
        }

        // Append the repositories to the repos slice
        repos = append(repos, pageRepos...)

        // Get the next page URL from the Link header
        nextPageURL = getNextPageURL(resp.Header.Get("Link"))

        // Wait for 1 minute before fetching the next page
        time.Sleep(1 * time.Minute)
    }

    return repos, nil
}

// Get the next page URL from the Link header
func getNextPageURL(linkHeader string) string {
    if linkHeader == "" {
        return ""
    }

    links := strings.Split(linkHeader, ",")
    for _, link := range links {
        parts := strings.Split(link, ";")
        if len(parts) == 2 && strings.TrimSpace(parts[1]) == `rel="next"` {
            return strings.TrimSpace(parts[0])[1 : len(parts[0])-2]
        }
    }

    return ""
}

// Get the total star count for a list of repositories
func getTotalStars(repos []Repository) int {
    totalStars := 0
    for _, repo := range repos {
        totalStars += repo.Stars
    }
    return totalStars
}

func main() {
    username := "xyproto"
    if len(os.Args) > 1 {
        username = os.Args[1]
    }

    var currentStars int

    // Fetch the list of repositories
    repos, err := fetchRepositories(username)
    if err != nil {
        panic(err)
    }

    fmt.Printf("got %d repos\n", len(repos))

    for {
        // Get the total star count
        totalStars := getTotalStars(repos)

        fmt.Printf("got %d total stars\n", totalStars)

        // If this is the first iteration, initialize currentStars
        if currentStars == 0 {
            currentStars = totalStars
        }

        // If the star count has increased, print "STAR!" and update currentStars
        if totalStars > currentStars {
            fmt.Println("STAR!")
            currentStars = totalStars
        }

        // Wait for 10 minutes before checking again
        time.Sleep(10 * time.Minute)
    }
}
