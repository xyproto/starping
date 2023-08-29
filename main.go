package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Repository represents a GitHub repository
type Repository struct {
	Name  string `json:"name"`
	Stars int    `json:"stargazers_count"`
}

// Fetch the list of repositories for a given GitHub user
func fetchRepositories(username string) ([]Repository, error) {
	// Fetch the list of repositories
	resp, err := http.Get("https://api.github.com/users/" + username + "/repos")
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
	var repos []Repository
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, err
	}

	return repos, nil
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

	for {
		// Fetch the list of repositories
		repos, err := fetchRepositories(username)
		if err != nil {
			panic(err)
		}

		// Get the total star count
		totalStars := getTotalStars(repos)

		// If this is the first iteration, initialize currentStars
		if currentStars == 0 {
			currentStars = totalStars
		}

		// If the star count has increased, print "STAR!" and update currentStars
		if totalStars > currentStars {
			fmt.Println("STAR!")
			currentStars = totalStars
		}

		// Wait for 1 minute before checking again
		time.Sleep(1 * time.Minute)
	}
}
