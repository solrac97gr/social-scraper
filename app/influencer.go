package app

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/solrac97gr/telegram-followers-checker/database"
	"github.com/solrac97gr/telegram-followers-checker/extractors/extractor"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
	ruregistration "github.com/solrac97gr/telegram-followers-checker/ru-registration"
)

// InfluencerApp orchestrates the components of the application
type InfluencerApp struct {
	influencersRepository database.InfluencerRepository
	fileManager           filemanager.FileManager
	extractors            []extractor.StatisticExtractor
}

// NewInfluencerApp creates a new App instance
func NewInfluencerApp(influencersRepository database.InfluencerRepository, fm filemanager.FileManager, extractors ...extractor.StatisticExtractor) *InfluencerApp {

	if influencersRepository == nil {
		log.Fatal("influencersRepository cannot be nil")
	}
	if fm == nil {
		log.Fatal("fileManager cannot be nil")
	}
	if len(extractors) == 0 {
		log.Fatal("At least one extractor must be provided")
	}

	return &InfluencerApp{
		influencersRepository: influencersRepository,
		fileManager:           fm,
		extractors:            extractors,
	}
}

// Run processes the input file and generates the output file
func (a *InfluencerApp) Run(inputFile string, outputFile string) [][]string {
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
		resp, err := a.influencersRepository.GetInfluencerAnalysisByLink(link)
		if err != nil {
			log.Printf("Error fetching analysis for %s: %v", link, err)
		}
		if resp != nil && err == nil {
			log.Printf("Link %s already processed, getting from database.", link)
			resultsList[i] = resp.ToExcelRow()
		} else { // Find appropriate extractor for this link
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
				// Store result directly at the correct position
				analysis := database.NewInfluencerAnalysis(
					info.ChannelName,        // ChannelName
					info.OriginalLink,       // Link
					info.Platform,           // Platform
					info.FollowersCount,     // FollowersCount
					info.RegistrationStatus, // RegistrationStatus
				)
				resultsList[i] = analysis.ToExcelRow()
				err := a.influencersRepository.SaveInfluencerAnalysis(analysis)
				if err != nil {
					log.Printf("Error saving analysis for %s: %v", info.OriginalLink, err)
					continue
				}
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

				// Store result at the correct position in the resultsList
				mutex.Lock()
				analysis := database.NewInfluencerAnalysis(
					currentInfo.ChannelName,        // ChannelName
					currentInfo.OriginalLink,       // Link
					currentInfo.Platform,           // Platform
					currentInfo.FollowersCount,     // FollowersCount
					currentInfo.RegistrationStatus, // RegistrationStatus
				)
				err := a.influencersRepository.SaveInfluencerAnalysis(analysis)
				if err != nil {
					log.Printf("Error saving analysis for %s: %v", info.OriginalLink, err)
				}
				resultsList[idx] = analysis.ToExcelRow()
				mutex.Unlock()
				// Avoid hitting rate limits
				time.Sleep(1 * time.Second)
			}(i, info, link)
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Append all results in order
	orderedResults = append(orderedResults, resultsList...)

	log.Printf("Processed %d links successfully. Preparing to save results...", len(links))

	log.Printf("orderedResults: %v", orderedResults)
	// Save results to output file
	a.fileManager.SaveResultsToExcel(orderedResults, outputFile)

	return orderedResults
}

func (a *InfluencerApp) GetAllInfluencerAnalysis(pageNum, limit int) (database.AllInfluencerAnalysis, error) {
	return a.influencersRepository.GetAllInfluencerAnalyses(pageNum, limit)
}

func (a *InfluencerApp) EstimateProcessingTime(inputFile string) (int, error) {
	return a.fileManager.EstimateProcessingTime(inputFile)
}
