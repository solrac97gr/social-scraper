package instagram

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/solrac97gr/telegram-followers-checker/extractors/extractor"
)

const (
	InstagramUsernameEnv = "INSTAGRAM_USERNAME"
	InstagramPasswordEnv = "INSTAGRAM_PASSWORD"
)

// InstagramExtractor is an implementation of StatisticExtractor for Instagram
type InstagramExtractor struct {
	name string
}

// NewInstagramExtractor creates a new InstagramExtractor instance
func NewInstagramExtractor() *InstagramExtractor {
	return &InstagramExtractor{
		name: "instagram",
	}
}

// Name returns the name of this extractor
func (ie *InstagramExtractor) Name() string {
	return ie.name
}

// CanHandle returns true if this extractor can handle the given link
func (ie *InstagramExtractor) CanHandle(link string) bool {
	return strings.Contains(link, "instagram.com/")
}

// Extract extracts channel information from the given link
func (ie *InstagramExtractor) Extract(link string) extractor.ChannelInfo {
	// Run the Node.js script using Puppeteer
	cmd := exec.Command(
		"node",
		"scripts/instagram.js",
		link,
		os.Getenv(InstagramUsernameEnv),
		os.Getenv(InstagramPasswordEnv),
	)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running Puppeteer script: %v", err)
		return extractor.ChannelInfo{
			ChannelName:    "Error",
			FollowersCount: "N/A",
			OriginalLink:   link,
		}
	}

	// Parse the JSON output from the Node.js script
	var result struct {
		ChannelName    string `json:"channelName"`
		FollowersCount string `json:"followersCount"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		log.Printf("Error parsing JSON output: %v", err)
		return extractor.ChannelInfo{
			ChannelName:    "Error",
			FollowersCount: "N/A",
			OriginalLink:   link,
		}
	}

	// Convert followers count to a number
	result.FollowersCount = convertFollowersCount(result.FollowersCount)

	return extractor.ChannelInfo{
		ChannelName:    result.ChannelName,
		FollowersCount: result.FollowersCount,
		OriginalLink:   link,
	}
}

// convertFollowersCount converts followers count from a string like "10K" to a number string like "10000"
func convertFollowersCount(followersText string) string {
	followersText = strings.ToUpper(followersText)
	if strings.Contains(followersText, "K") {
		followersText = strings.ReplaceAll(followersText, "K", "000")
	} else if strings.Contains(followersText, "M") {
		followersText = strings.ReplaceAll(followersText, "M", "000000")
	}
	return followersText
}
