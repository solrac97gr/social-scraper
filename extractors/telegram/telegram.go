package telegram

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"

	"github.com/solrac97gr/telegram-followers-checker/extractors/extractor"
)

// TelegramExtractor is an implementation of StatisticExtractor for Telegram
type TelegramExtractor struct {
	name string
}

// NewTelegramExtractor creates a new TelegramExtractor instance
func NewTelegramExtractor() *TelegramExtractor {
	return &TelegramExtractor{
		name: "telegram",
	}
}

// Name returns the name of this extractor
func (te *TelegramExtractor) Name() string {
	return te.name
}

// CanHandle returns true if this extractor can handle the given link
func (te *TelegramExtractor) CanHandle(link string) bool {
	return strings.Contains(link, "t.me/") || strings.Contains(link, "telegram.me/")
}

// Extract extracts channel information from the given link using Puppeteer script
func (te *TelegramExtractor) Extract(link string) extractor.ChannelInfo {
	// Format the link to ensure it's accessible via http
	if !strings.HasPrefix(link, "http") {
		link = "https://" + link
	}

	// Execute the Telegram Puppeteer script
	cmd := exec.Command("node", "scripts/telegram.js", link)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing Telegram script for %s: %v", link, err)
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
		log.Printf("Error parsing Telegram script output for %s: %v", link, err)
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
