package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (repo *MongoRepository) DeleteExpiredAnalyses() error {
	ctx := context.Background()
	collection := repo.client.Database(DatabaseName).Collection(CollectionName)
	filter := bson.M{"expiration_date": bson.M{"$lt": time.Now()}}
	_, err := collection.DeleteMany(ctx, filter)
	return err
}

func (repo *MongoRepository) GetAllInfluencerAnalyses(page int, limit int) (AllInfluencerAnalysis, error) {
	ctx := context.Background()
	collection := repo.client.Database(DatabaseName).Collection(CollectionName)

	skip := (page - 1) * limit
	filter := bson.M{"expiration_date": bson.M{"$gt": time.Now()}}
	cursor, err := collection.Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		return AllInfluencerAnalysis{}, err
	}
	defer cursor.Close(ctx)

	var analyses []*InfluencerAnalysis
	for cursor.Next(ctx) {
		var analysis InfluencerAnalysis
		if err := cursor.Decode(&analysis); err != nil {
			return AllInfluencerAnalysis{}, err
		}
		analyses = append(analyses, &analysis)
	}

	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return AllInfluencerAnalysis{}, err
	}

	return AllInfluencerAnalysis{
		TotalCount: totalCount,
		Analyses:   analyses,
		Pagination: struct {
			Page  int64 `json:"page" bson:"page"`
			Limit int64 `json:"limit" bson:"limit"`
		}{
			Page:  int64(page),
			Limit: int64(limit),
		},
	}, nil
}
