package telegram

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

// Extract extracts channel information from the given link
func (te *TelegramExtractor) Extract(link string) extractor.ChannelInfo {
	// Format the link to ensure it's accessible via http
	if !strings.HasPrefix(link, "http") {
		link = "https://" + link
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	resp, err := client.Get(link)
	if err != nil {
		log.Printf("Error fetching %s: %v", link, err)
		return extractor.ChannelInfo{
			ChannelName:    "Error",
			FollowersCount: "N/A",
			OriginalLink:   link,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Status code error: %d %s", resp.StatusCode, resp.Status)
		return extractor.ChannelInfo{
			ChannelName:    "Error",
			FollowersCount: "N/A",
			OriginalLink:   link,
		}
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		return extractor.ChannelInfo{
			ChannelName:    "Error",
			FollowersCount: "N/A",
			OriginalLink:   link,
		}
	}

	// Extract channel name
	channelName := "Unknown"
	doc.Find("div.tgme_page_title").Each(func(i int, s *goquery.Selection) {
		channelName = strings.TrimSpace(s.Text())
	})

	// Extract followers count
	followersText := "0"
	doc.Find("div.tgme_page_extra").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, "subscriber") || strings.Contains(text, "member") || strings.Contains(text, "follower") {
			re := regexp.MustCompile(`[\d\s]+`)
			matches := re.FindString(text)
			if matches != "" {
				// Remove spaces and convert to number
				followersText = strings.ReplaceAll(matches, " ", "")
			}
		}
	})

	return extractor.ChannelInfo{
		ChannelName:    channelName,
		FollowersCount: followersText,
		OriginalLink:   link,
	}
}
