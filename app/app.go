package app

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/solrac97gr/telegram-followers-checker/extractors/extractor"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
	ruregistration "github.com/solrac97gr/telegram-followers-checker/ru-registration"
)

// App orchestrates the components of the application
type App struct {
	fileManager filemanager.FileManager
	extractors  []extractor.StatisticExtractor
}

// NewApp creates a new App instance
func NewApp(fm filemanager.FileManager, extractors ...extractor.StatisticExtractor) *App {
	return &App{
		fileManager: fm,
		extractors:  extractors,
	}
}

// Run processes the input file and generates the output file
func (a *App) Run(inputFile string, outputFile string) [][]string {
	// Read links from input file
	links := a.fileManager.ReadLinksFromExcel(inputFile)

	// Create a slice to store results in order
	orderedResults := make([][]string, 0, len(links)+1)
	// Add header row
	orderedResults = append(orderedResults, []string{"Channel Name", "Followers Count", "Original Link", "Platform", "Registration Status"})

	// Create a channel to collect results
	resultsChan := make(chan []string, len(links))

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a semaphore to limit concurrent registration checks
	semaphore := make(chan struct{}, 10)

	// Process each link concurrently
	for i, link := range links {
		// Find appropriate extractor for this link
		var info extractor.ChannelInfo
		for _, e := range a.extractors {
			if e.CanHandle(link) {
				info = e.Extract(link)
				info.Platform = e.Name()
				break
			}
		}

		// If no extractor found or extraction failed, use defaults
		if info.ChannelName == "" {
			info = extractor.ChannelInfo{
				ChannelName:    "Unknown",
				FollowersCount: "0",
				OriginalLink:   link,
				Platform:       "Unknown",
			}
		}

		// Skip registration status check if platform is Instagram or followers count is < 10000
		followersCount, err := strconv.Atoi(info.FollowersCount)
		if info.Platform == "Instagram" || (err == nil && followersCount < 10000) {
			info.RegistrationStatus = "not applicable ⚪"
			resultsChan <- []string{info.ChannelName, info.FollowersCount, info.OriginalLink, info.Platform, info.RegistrationStatus}
			continue
		}

		// Add to WaitGroup only for links that will be processed
		wg.Add(1)
		go func(i int, link string) {
			defer wg.Done()

			// Define isRegistered channel
			isRegistered := make(chan bool)

			go func() {
				semaphore <- struct{}{} // Acquire semaphore
				isRegistered <- ruregistration.CheckRegistrationStatus(link, semaphore)
				close(isRegistered)
			}()

			// Collect the result
			info.IsRegistered = <-isRegistered
			if info.IsRegistered {
				info.RegistrationStatus = "registered 🟢"
			} else {
				info.RegistrationStatus = "not registered 🔴"
			}

			// Send the result to the channel
			resultsChan <- []string{info.ChannelName, info.FollowersCount, info.OriginalLink, info.Platform, info.RegistrationStatus}

			// Avoid hitting rate limits
			time.Sleep(1 * time.Second)
		}(i, link)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(resultsChan)

	// Collect results from the channel
	for result := range resultsChan {
		orderedResults = append(orderedResults, result)
	}

	// Save results to output file
	a.fileManager.SaveResultsToExcel(orderedResults, outputFile)

	fmt.Printf("\nSuccess! Results saved to %s\n", outputFile)

	return orderedResults
}
