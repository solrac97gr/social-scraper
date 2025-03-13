package app

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/solrac97gr/telegram-followers-checker/extractors/extractor"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
	ruregistration "github.com/solrac97gr/telegram-followers-checker/integrations/ru-registration"
	"github.com/solrac97gr/telegram-followers-checker/integrations/tgstats"
	"github.com/solrac97gr/telegram-followers-checker/integrations/tgstats/config"
	"github.com/solrac97gr/telegram-followers-checker/integrations/tgstats/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// App orchestrates the components of the application
type App struct {
	fileManager filemanager.FileManager
	extractors  []extractor.StatisticExtractor
	repo        *repository.MongoRepository
	config      *config.TGStatsConfig
}

// NewApp creates a new App instance
func NewApp(fm filemanager.FileManager, repo *repository.MongoRepository, config *config.TGStatsConfig, extractors ...extractor.StatisticExtractor) *App {
	return &App{
		fileManager: fm,
		extractors:  extractors,
		repo:        repo,
		config:      config,
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

	// Create a slice to store results at the correct index
	resultsList := make([][]string, len(links))

	// Create a mutex to protect the resultsList during concurrent writes
	var mutex sync.Mutex

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a semaphore to limit concurrent registration checks
	semaphore := make(chan struct{}, 10)

	// Process each link concurrently
	for i, link := range links {
		// Find appropriate extractor for this link
		var info extractor.ChannelInfo
		// Check cache before performing extraction
		cachedInfo, err := a.repo.GetChannelInfo(link)
		if err == nil && cachedInfo.ExpirationTime.After(time.Now()) {
			info = *cachedInfo
		} else {
			// Perform extraction if not in cache or expired
			for _, e := range a.extractors {
				if e.CanHandle(link) {
					info = e.Extract(link)
					info.Platform = e.Name()
					info.ExpirationTime = time.Now().Add(72 * time.Hour)
					a.repo.SaveChannelInfo(&info)
					break
				}
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
		if info.Platform == "instagram" || (err == nil && followersCount < 10000) {
			info.RegistrationStatus = "not applicable ⚪"
			// Store result directly at the correct position
			resultsList[i] = []string{info.ChannelName, info.FollowersCount, info.OriginalLink, info.Platform, info.RegistrationStatus, fmt.Sprintf("%.2f", info.AvgPostReach), fmt.Sprintf("%.2f", info.ERPercent)}
			continue
		}

		// Add to WaitGroup only for links that will be processed
		wg.Add(1)
		go func(idx int, currentInfo extractor.ChannelInfo, linkUrl string) {
			defer wg.Done()

			// Define isRegistered channel
			isRegistered := make(chan bool)

			go func() {
				semaphore <- struct{}{} // Acquire semaphore
				isRegistered <- ruregistration.CheckRegistrationStatus(linkUrl, semaphore)
				close(isRegistered)
			}()

			// Collect the result
			currentInfo.IsRegistered = <-isRegistered
			if currentInfo.IsRegistered {
				currentInfo.RegistrationStatus = "registered 🟢"
			} else {
				currentInfo.RegistrationStatus = "not registered 🔴"
			}

			// Get TGStats data if the platform is Telegram
			if currentInfo.Platform == "telegram" {
				tgStatsResult, err := tgstats.GetTGStats(linkUrl, a.repo, a.config)
				if err == nil {
					currentInfo.AvgPostReach = tgStatsResult.AvgPostReach
					currentInfo.ERPercent = tgStatsResult.ERPercent
				}
			}

			// Store result at the correct position in the resultsList
			mutex.Lock()
			resultsList[idx] = []string{currentInfo.ChannelName, currentInfo.FollowersCount, currentInfo.OriginalLink, currentInfo.Platform, currentInfo.RegistrationStatus, fmt.Sprintf("%.2f", currentInfo.AvgPostReach), fmt.Sprintf("%.2f", currentInfo.ERPercent)}
			mutex.Unlock()

			// Avoid hitting rate limits
			time.Sleep(1 * time.Second)
		}(i, info, link)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Append all results in order
	orderedResults = append(orderedResults, resultsList...)

	// Save results to output file
	a.fileManager.SaveResultsToExcel(orderedResults, outputFile)

	fmt.Printf("\nSuccess! Results saved to %s\n", outputFile)

	return orderedResults
}

func InitMongoClient(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return client, nil
}
