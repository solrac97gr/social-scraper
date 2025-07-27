package rutube

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"

	"github.com/solrac97gr/telegram-followers-checker/extractors/extractor"
)

// RutubeExtractor is an implementation of StatisticExtractor for Rutube
type RutubeExtractor struct {
	name string
}

// NewRutubeExtractor creates a new RutubeExtractor instance
func NewRutubeExtractor() *RutubeExtractor {
	return &RutubeExtractor{
		name: "rutube",
	}
}

// Name returns the name of this extractor
func (re *RutubeExtractor) Name() string {
	return re.name
}

// CanHandle returns true if this extractor can handle the given link
func (re *RutubeExtractor) CanHandle(link string) bool {
	canHandle := strings.Contains(link, "rutube.ru/")
	return canHandle
}

// Extract extracts channel information from the given link using Puppeteer script
func (re *RutubeExtractor) Extract(link string) extractor.ChannelInfo {
	// Format the link to ensure it's accessible via http
	if !strings.HasPrefix(link, "http") {
		link = "https://" + link
	}

	// Execute the Rutube Puppeteer script
	cmd := exec.Command("node", "scripts/rutube.js", link)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing Rutube script for %s: %v", link, err)
		return extractor.ChannelInfo{
			ChannelName:    "Error",
			FollowersCount: "N/A",
			OriginalLink:   link,
		}
	}

	// Parse the JSON response from the script
	var result struct {
		ChannelName    string `json:"channelName"`
		FollowersCount string `json:"followersCount"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		log.Printf("Error parsing Rutube script output for %s: %v", link, err)
		return extractor.ChannelInfo{
			ChannelName:    "Error",
			FollowersCount: "N/A",
			OriginalLink:   link,
		}
	}

	return extractor.ChannelInfo{
		ChannelName:    result.ChannelName,
		FollowersCount: result.FollowersCount,
		OriginalLink:   link,
	}
}
