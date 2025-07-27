package tiktok

import (
	"encoding/json"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/solrac97gr/telegram-followers-checker/extractors/extractor"
)

// TikTokExtractor is an implementation of StatisticExtractor for TikTok
type TikTokExtractor struct {
	name string
}

// NewTikTokExtractor creates a new TikTokExtractor instance
func NewTikTokExtractor() *TikTokExtractor {
	return &TikTokExtractor{
		name: "tiktok",
	}
}

// Name returns the name of this extractor
func (te *TikTokExtractor) Name() string {
	return te.name
}

// CanHandle returns true if this extractor can handle the given link
func (te *TikTokExtractor) CanHandle(link string) bool {
	return strings.Contains(link, "tiktok.com/")
}

// Extract extracts channel information from the given link
func (te *TikTokExtractor) Extract(link string) extractor.ChannelInfo {
	// Ensure link starts with https://
	if !strings.HasPrefix(link, "http") {
		link = "https://" + link
	} else if strings.HasPrefix(link, "http://") {
		link = "https://" + link[len("http://"):]
	}

	// Remove any subdomain from tiktok.com and ensure canonical format
	re := regexp.MustCompile(`https://[^/]*tiktok\.com`)
	modifiedLink := re.ReplaceAllString(link, `https://www.tiktok.com`)

	println("Modified TikTok link:", modifiedLink)

	// Run the Node.js script using Puppeteer
	cmd := exec.Command("node", "scripts/tiktok.js", modifiedLink)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running TikTok Puppeteer script: %v", err)
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
		log.Printf("Output: %s", string(output))
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

// convertFollowersCount converts followers count from formats like "20K", "1.5M", "5000" to a number string
func convertFollowersCount(followersText string) string {
	// Clean the input
	followersText = strings.TrimSpace(followersText)
	followersText = strings.ToUpper(followersText)

	// Remove any non-alphanumeric characters except dots and common suffixes
	re := regexp.MustCompile(`[^\d\.KM]`)
	followersText = re.ReplaceAllString(followersText, "")

	// Handle different formats
	if strings.Contains(followersText, "K") {
		// Handle formats like "20K" or "20.5K"
		numStr := strings.ReplaceAll(followersText, "K", "")
		if num, err := strconv.ParseFloat(numStr, 64); err == nil {
			return strconv.Itoa(int(num * 1000))
		}
	} else if strings.Contains(followersText, "M") {
		// Handle formats like "1M" or "1.5M"
		numStr := strings.ReplaceAll(followersText, "M", "")
		if num, err := strconv.ParseFloat(numStr, 64); err == nil {
			return strconv.Itoa(int(num * 1000000))
		}
	} else {
		// Handle plain numbers, remove any commas
		followersText = strings.ReplaceAll(followersText, ",", "")
		if num, err := strconv.Atoi(followersText); err == nil {
			return strconv.Itoa(num)
		}
	}

	// If conversion fails, return the original text
	return followersText
}
