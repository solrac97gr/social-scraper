package app

import (
	"fmt"
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

	// Process each link
	for _, link := range links {
		fmt.Printf("Processing: %s\n", link)

		// Find appropriate extractor for this link
		var info extractor.ChannelInfo
		for _, e := range a.extractors {
			if e.CanHandle(link) {
				info = e.Extract(link)
				fmt.Printf("Using %s extractor\n", e.Name())
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

		// Add to results
		results = append(results, []string{info.ChannelName, info.FollowersCount, info.OriginalLink})

		// Avoid hitting rate limits
		time.Sleep(1 * time.Second)
	}

	// Save results to output file
	a.fileManager.SaveResultsToExcel(results, outputFile)

	fmt.Printf("\nSuccess! Results saved to %s\n", outputFile)
}
