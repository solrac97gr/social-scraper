package repository

import (
	"context"
	"time"

	"github.com/solrac97gr/telegram-followers-checker/extractors/extractor"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongoRepository defines the structure for the MongoDB repository
type MongoRepository struct {
	client   *mongo.Client
	database string
}

// NewMongoRepository creates a new instance of MongoRepository
func NewMongoRepository(client *mongo.Client, database string) *MongoRepository {
	return &MongoRepository{
		client:   client,
		database: database,
	}
}

// ChannelData represents the data structure for channel information
// Updated to include all fields from the API response

type ChannelData struct {
	Channel                 string    `bson:"channel"`
	AvgPostReach            float32   `bson:"avg_post_reach"`
	ERPercent               float32   `bson:"er_percent"`
	ID                      int       `bson:"id"`
	Title                   string    `bson:"title"`
	Username                string    `bson:"username"`
	PeerType                string    `bson:"peer_type"`
	ParticipantsCount       int       `bson:"participants_count"`
	AdvPostReach12h         float32   `bson:"adv_post_reach_12h"`
	AdvPostReach24h         float32   `bson:"adv_post_reach_24h"`
	AdvPostReach48h         float32   `bson:"adv_post_reach_48h"`
	ErrPercent              float32   `bson:"err_percent"`
	Err24Percent            float32   `bson:"err24_percent"`
	DailyReach              float32   `bson:"daily_reach"`
	CiIndex                 float32   `bson:"ci_index"`
	MentionsCount           int       `bson:"mentions_count"`
	ForwardsCount           int       `bson:"forwards_count"`
	MentioningChannelsCount int       `bson:"mentioning_channels_count"`
	PostsCount              int       `bson:"posts_count"`
	ExpirationTime          time.Time `bson:"expiration_time"`
}

// RequestCount represents the data structure for request count information
type RequestCount struct {
	RequestAmount    int    `bson:"request_amount"`
	RequestNotCached int    `bson:"request_not_cached"`
	DatePeriod       string `bson:"date_period"`
}

// GetChannelData retrieves the channel data from the database
func (r *MongoRepository) GetChannelData(channel string) (*ChannelData, error) {
	collection := r.client.Database(r.database).Collection("channels")
	var data ChannelData
	err := collection.FindOne(context.TODO(), bson.M{"channel": channel}).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// SaveChannelData saves the channel data to the database
func (r *MongoRepository) SaveChannelData(data *ChannelData) error {
	collection := r.client.Database(r.database).Collection("channels")
	_, err := collection.InsertOne(context.TODO(), data)
	return err
}

// GetRequestCount retrieves the request count from the database
func (r *MongoRepository) GetRequestCount(datePeriod string) (*RequestCount, error) {
	collection := r.client.Database(r.database).Collection("requests")
	var count RequestCount
	err := collection.FindOne(context.TODO(), bson.M{"date_period": datePeriod}).Decode(&count)
	if err != nil {
		return nil, err
	}
	return &count, nil
}

// SaveRequestCount saves the request count to the database
func (r *MongoRepository) SaveRequestCount(count *RequestCount) error {
	collection := r.client.Database(r.database).Collection("requests")
	_, err := collection.InsertOne(context.TODO(), count)
	return err
}

// SaveChannelInfo saves the channel info to the database
func (r *MongoRepository) SaveChannelInfo(info *extractor.ChannelInfo) error {
	collection := r.client.Database(r.database).Collection("channel_info")
	_, err := collection.InsertOne(context.TODO(), info)
	return err
}

// GetChannelInfo retrieves the channel info from the database
func (r *MongoRepository) GetChannelInfo(OriginalLink string) (*extractor.ChannelInfo, error) {
	collection := r.client.Database(r.database).Collection("channel_info")
	var info extractor.ChannelInfo
	err := collection.FindOne(context.TODO(), bson.M{"originallink": OriginalLink}).Decode(&info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
