package vk

import (
	"encoding/json"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/solrac97gr/telegram-followers-checker/extractors/extractor"
)

// VKExtractor is an implementation of StatisticExtractor for VK
type VKExtractor struct {
	name string
}

// NewVKExtractor creates a new VKExtractor instance
func NewVKExtractor() *VKExtractor {
	return &VKExtractor{
		name: "vk",
	}
}

// Name returns the name of this extractor
func (ve *VKExtractor) Name() string {
	return ve.name
}

// CanHandle returns true if this extractor can handle the given link
func (ve *VKExtractor) CanHandle(link string) bool {
	return strings.Contains(link, "vk.com/")
}

// Extract extracts channel information from the given link
func (ve *VKExtractor) Extract(link string) extractor.ChannelInfo {
	// Ensure link starts with https://
	if !strings.HasPrefix(link, "http") {
		link = "https://" + link
	} else if strings.HasPrefix(link, "http://") {
		link = "https://" + link[len("http://"):]
	}

	// Remove any subdomain from vk.com
	re := regexp.MustCompile(`https://[^/]*vk\.com`)
	modifiedLink := re.ReplaceAllString(link, `https://vk.com`)

	println("Modified VK link:", modifiedLink)

	// Run the Node.js script using Puppeteer
	cmd := exec.Command("node", "scripts/puppeteer_scraper.js", modifiedLink)
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
		FollowersCount string `json:"followersText"`
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

// convertFollowersCount converts followers count from a string like "10K" or "3,700" to a number string like "10000" or "3700"
func convertFollowersCount(followersText string) string {
	followersText = strings.ToUpper(followersText)
	followersText = strings.ReplaceAll(followersText, ",", "")
	if strings.Contains(followersText, "K") {
		followersText = strings.ReplaceAll(followersText, "K", "000")
	} else if strings.Contains(followersText, "M") {
		followersText = strings.ReplaceAll(followersText, "M", "000000")
	}
	return followersText
}
