package extractor

// ChannelInfo holds information about a channel
type ChannelInfo struct {
	ChannelName    string
	FollowersCount string
	OriginalLink   string
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
