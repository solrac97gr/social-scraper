// Description: This package defines the interface for extracting statistics from channels using scrapers.
package extractor

import "time"

// ChannelInfo holds information about a channel
type ChannelInfo struct {
	ChannelName        string
	FollowersCount     string
	OriginalLink       string
	Platform           string
	IsRegistered       bool
	RegistrationStatus string
	AvgPostReach       float32
	ERPercent          float32
	ExpirationTime     time.Time `bson:"expiration_time"`
}

// StatisticExtractor defines the interface for extracting statistics from channels
type StatisticExtractor interface {
	// CanHandle returns true if this extractor can handle the given link
	CanHandle(link string) bool

	// Extract extracts channel information from the given link
	Extract(link string) ChannelInfo

	// Name returns the name of this extractor
	Name() string
}
