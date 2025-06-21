package rutube

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

// Extract extracts channel information from the given link
func (re *RutubeExtractor) Extract(link string) extractor.ChannelInfo {
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

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
	channelName := doc.Find("h1.wdp-feed-banner-module__wdp-feed-banner__title-text").AttrOr("title", "")
	channelName = strings.TrimSpace(channelName)

	// Extract followers count
	followersText := doc.Find(".wdp-feed-banner-module__wdp-feed-banner__title p").Text()
	followersText = strings.TrimSpace(followersText)

	// Clean up the followers string to remove non-digit characters
	regex := regexp.MustCompile(`\D`)
	followersText = regex.ReplaceAllString(followersText, "")

	if channelName == "" || followersText == "" {
		return extractor.ChannelInfo{
			ChannelName:    "Error",
			FollowersCount: "N/A",
			OriginalLink:   link,
		}
	}

	return extractor.ChannelInfo{
		ChannelName:    channelName,
		FollowersCount: followersText,
		OriginalLink:   link,
	}
}
