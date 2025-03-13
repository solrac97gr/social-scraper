package tgstats

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/solrac97gr/telegram-followers-checker/integrations/tgstats/config"
	"github.com/solrac97gr/telegram-followers-checker/integrations/tgstats/repository"
)

type TGStatsResult struct {
	AvgPostReach float32
	ERPercent    float32
}

type TGStatAPIResponse struct {
	Status   string `json:"status"`
	Response struct {
		AvgPostReach            float32 `json:"avg_post_reach"`
		ERPercent               float32 `json:"er_percent"`
		ID                      int     `json:"id"`
		Username                string  `json:"username"`
		Title                   string  `json:"title"`
		PeerType                string  `json:"peer_type"`
		ParticipantsCount       int     `json:"participants_count"`
		AdvPostReach12h         float32 `json:"adv_post_reach_12h"`
		AdvPostReach24h         float32 `json:"adv_post_reach_24h"`
		AdvPostReach48h         float32 `json:"adv_post_reach_48h"`
		ErrPercent              float32 `json:"err_percent"`
		Err24Percent            float32 `json:"err_24_percent"`
		DailyReach              float32 `json:"daily_reach"`
		CiIndex                 float32 `json:"ci_index"`
		MentionsCount           int     `json:"mentions_count"`
		ForwardsCount           int     `json:"forwards_count"`
		MentioningChannelsCount int     `json:"mentioning_channels_count"`
		PostsCount              int     `json:"posts_count"`
	} `json:"response"`
}

func GetTGStats(channel string, repo *repository.MongoRepository, config *config.TGStatsConfig) (*TGStatsResult, error) {
	log.Printf("Checking cache for channel: %s", channel)
	data, err := repo.GetChannelData(channel)
	if err == nil && data.ExpirationTime.After(time.Now()) {
		log.Printf("Cache hit for channel: %s", channel)
		return &TGStatsResult{
			AvgPostReach: data.AvgPostReach,
			ERPercent:    data.ERPercent,
		}, nil
	}
	log.Printf("Cache miss or expired data for channel: %s", channel)

	// Make API request
	url := fmt.Sprintf("%s?token=%s&channelId=%s", config.URL, config.Token, channel)
	log.Printf("Making API request to URL: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error making API request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading API response body: %v", err)
		return nil, err
	}
	log.Printf("API response body: %s", string(body))

	var apiResponse TGStatAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Error unmarshalling API response: %v", err)
		return nil, err
	}

	if apiResponse.Status != "ok" {
		log.Printf("API response status not ok: %s", apiResponse.Status)
		return nil, fmt.Errorf("API response status: %s", apiResponse.Status)
	}

	// Save data to cache
	log.Printf("Saving data to cache for channel: %s", channel)
	newData := &repository.ChannelData{
		Channel:                 channel,
		AvgPostReach:            apiResponse.Response.AvgPostReach,
		ERPercent:               apiResponse.Response.ERPercent,
		ID:                      apiResponse.Response.ID,
		Title:                   apiResponse.Response.Title,
		Username:                apiResponse.Response.Username,
		PeerType:                apiResponse.Response.PeerType,
		ParticipantsCount:       apiResponse.Response.ParticipantsCount,
		AdvPostReach12h:         apiResponse.Response.AdvPostReach12h,
		AdvPostReach24h:         apiResponse.Response.AdvPostReach24h,
		AdvPostReach48h:         apiResponse.Response.AdvPostReach48h,
		ErrPercent:              apiResponse.Response.ErrPercent,
		Err24Percent:            apiResponse.Response.Err24Percent,
		DailyReach:              apiResponse.Response.DailyReach,
		CiIndex:                 apiResponse.Response.CiIndex,
		MentionsCount:           apiResponse.Response.MentionsCount,
		ForwardsCount:           apiResponse.Response.ForwardsCount,
		MentioningChannelsCount: apiResponse.Response.MentioningChannelsCount,
		PostsCount:              apiResponse.Response.PostsCount,
		ExpirationTime:          time.Now().AddDate(0, 1, 0), // 1 month expiration
	}
	repo.SaveChannelData(newData)

	log.Printf("Returning TGStatsResult for channel: %s", channel)
	return &TGStatsResult{
		AvgPostReach: apiResponse.Response.AvgPostReach,
		ERPercent:    apiResponse.Response.ERPercent,
	}, nil
}
