package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/solrac97gr/telegram-followers-checker/extractors/extractor"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
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
func (a *App) Run(inputFile string, outputFile string) {
	// Read links from input file
	links := a.fileManager.ReadLinksFromExcel(inputFile)

	results := make([][]string, 0, len(links)+1)
	// Add header row
	results = append(results, []string{"Channel Name", "Followers Count", "Original Link"})

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(len(links))

	// Create a channel to collect the results
	resultChan := make(chan []string, len(links))

	// Process each link concurrently
	for i, link := range links {
		go func(i int, link string) {
			defer wg.Done()
			fmt.Printf("Processing: %s\n", link)

			// Find appropriate extractor for this link
			var info extractor.ChannelInfo
			for _, e := range a.extractors {
				if e.CanHandle(link) {
					info = e.Extract(link)
					break
				}
			}

			// If no extractor found or extraction failed, use defaults
			if info.ChannelName == "" {
				info = extractor.ChannelInfo{
					ChannelName:    "Unknown",
					FollowersCount: "0",
					OriginalLink:   link,
				}
			}

			// Send the result to the channel
			resultChan <- []string{info.ChannelName, info.FollowersCount, info.OriginalLink}

			// Avoid hitting rate limits
			time.Sleep(1 * time.Second)
		}(i, link)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(resultChan)

	// Collect results from the channel
	for result := range resultChan {
		results = append(results, result)
	}

	// Save results to output file
	a.fileManager.SaveResultsToExcel(results, outputFile)

	fmt.Printf("\nSuccess! Results saved to %s\n", outputFile)
}
