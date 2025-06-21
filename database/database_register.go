package database

import (
	"strconv"
	"time"
)

type Status string

const (
	Registered    Status = "registered"
	NotRegistered Status = "not_registered"
	NotApply      Status = "not_apply"
)

type InfluencerAnalysis struct {
	ID                 string    `json:"id" bson:"_id,omitempty"`
	UserID             string    `json:"user_id" bson:"user_id"` // ID of the user who created the analysis
	ChannelName        string    `json:"channel_name" bson:"channel_name"`
	FollowersCount     int       `json:"followers_count" bson:"followers_count"`
	Link               string    `json:"link" bson:"link"`
	Platform           string    `json:"platform" bson:"platform"`
	RegistrationStatus Status    `json:"registration_status" bson:"registration_status"`
	ExpirationDate     time.Time `json:"expiration_date" bson:"expiration_date"`
	CreatedAt          time.Time `json:"created_at" bson:"created_at"`
}

func NewInfluencerAnalysis(userID, channelName, link, platform string, followersCount string, registrationStatus string) *InfluencerAnalysis {
	var followersCountInt int
	followersCountInt, err := strconv.Atoi(followersCount)
	if err != nil {
		followersCountInt = 0 // Default to 0 if conversion fails
	}

	return &InfluencerAnalysis{
		UserID:         userID,
		ChannelName:    channelName,
		FollowersCount: followersCountInt,
		Link:           link,
		Platform:       platform,
		RegistrationStatus: func(rs string) Status {
			switch rs {
			case "registered 🟢":
				return Registered
			case "not registered 🔴":
				return NotRegistered
			default:
				return NotApply // Default to NotApply if not registered or not applicable
			}
		}(registrationStatus),
		ExpirationDate: time.Now().Add(30 * 24 * time.Hour), // Default expiration date set to 15 days from now
		CreatedAt:      time.Now(),
	}
}

func (dr *InfluencerAnalysis) ToExcelRow() []string {
	return []string{
		dr.ChannelName,
		strconv.Itoa(dr.FollowersCount),
		dr.Link,
		dr.Platform,
		dr.RegistrationStatus.String(),
	}
}

func (s Status) String() string {
	switch s {
	case Registered:
		return "registered 🟢"
	case NotRegistered:
		return "not registered 🔴"
	case NotApply:
		return "not applicable ⚪"
	default:
		return "unknown status"
	}
}
