package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DatabaseName   = "influencer-db"
	CollectionName = "influencer-analysis"
)

type MongoRepository struct {
	client *mongo.Client
}

func NewMongoRepository(client *mongo.Client) (*MongoRepository, error) {
	if client == nil {
		return nil, mongo.ErrClientDisconnected
	}

	return &MongoRepository{
		client: client,
	}, nil
}

func (repo *MongoRepository) save(ctx context.Context, data interface{}) error {
	collection := repo.client.Database(DatabaseName).Collection(CollectionName)
	_, err := collection.InsertOne(ctx, data)
	return err
}

func (repo *MongoRepository) findOne(ctx context.Context, filter interface{}) (*InfluencerAnalysis, error) {
	collection := repo.client.Database(DatabaseName).Collection(CollectionName)
	result := collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var data InfluencerAnalysis
	err := result.Decode(&data)
	return &data, err
}

func (repo *MongoRepository) SaveInfluencerAnalysis(influencer *InfluencerAnalysis) error {
	ctx := context.Background()
	return repo.save(ctx, influencer)

}

func (repo *MongoRepository) GetInfluencerAnalysisByLink(link string) (*InfluencerAnalysis, error) {
	ctx := context.Background()
	filter := bson.M{"link": link, "expiration_date": bson.M{"$gt": time.Now()}}
	analysis, err := repo.findOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return analysis, nil
}
