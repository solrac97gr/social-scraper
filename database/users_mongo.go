package database

import (
	"context"
	"errors"
	"time"

	"github.com/solrac97gr/telegram-followers-checker/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserAlreadyExist = errors.New("user already exists with this email")
)

const (
	UserDatabaseName          = "influencer-db"
	UserCollectionName        = "users"
	UserProfileCollectionName = "user-profiles"
	UserTokenCollectionName   = "user-tokens"
)

type UserMongoRepository struct {
	client *mongo.Client
	config *config.Config
}

var _ UserRepository = (*UserMongoRepository)(nil)

func NewUserMongoRepository(client *mongo.Client, config *config.Config) (*UserMongoRepository, error) {
	if client == nil {
		return nil, mongo.ErrClientDisconnected
	}

	return &UserMongoRepository{
		client: client,
		config: config,
	}, nil
}

// GetUserByID implements UserRepository.
func (u *UserMongoRepository) GetUserByID(userID string) (*User, error) {
	ctx := context.Background()
	collection := u.client.Database(u.config.UsersDBName).Collection(UserCollectionName)
	filter := bson.M{"_id": userID}
	result := collection.FindOne(ctx, filter)

	if result.Err() != nil {
		return nil, result.Err()
	}

	var user User
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail implements UserRepository.
func (u *UserMongoRepository) GetUserByEmail(email string) (*User, error) {
	ctx := context.Background()
	collection := u.client.Database(u.config.UsersDBName).Collection(UserCollectionName)
	filter := bson.M{"email": email}
	result := collection.FindOne(ctx, filter)

	if result.Err() != nil {
		return nil, result.Err()
	}

	var user User
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserTokenByUserID implements UserRepository.
func (u *UserMongoRepository) GetUserTokenByUserID(userID string) (*UserToken, error) {
	ctx := context.Background()
	collection := u.client.Database(u.config.UsersDBName).Collection(UserTokenCollectionName)
	filter := bson.M{"user_id": userID}
	result := collection.FindOne(ctx, filter)

	if result.Err() != nil {
		return nil, result.Err()
	}

	var token UserToken
	err := result.Decode(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// SaveUser implements UserRepository.
func (u *UserMongoRepository) SaveUser(user *User) (string, error) {
	ctx := context.Background()
	collection := u.client.Database(u.config.UsersDBName).Collection(UserCollectionName)
	// validate if the user already exists using the email field since it is unique
	filter := bson.M{"email": user.Email}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "", ErrUserAlreadyExist
	}
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// SaveUserProfile implements UserRepository.
func (u *UserMongoRepository) SaveUserProfile(profile *UserProfile) error {
	ctx := context.Background()
	collection := u.client.Database(u.config.UsersDBName).Collection(UserProfileCollectionName)
	_, err := collection.InsertOne(ctx, profile)
	return err
}

// SaveUserToken implements UserRepository.
func (u *UserMongoRepository) SaveUserToken(token *UserToken) error {
	ctx := context.Background()
	collection := u.client.Database(u.config.UsersDBName).Collection(UserTokenCollectionName)
	_, err := collection.InsertOne(ctx, token)
	return err
}

// UpdateUser implements UserRepository.
func (u *UserMongoRepository) UpdateUser(user *User) error {
	ctx := context.Background()
	collection := u.client.Database(u.config.UsersDBName).Collection(UserCollectionName)
	filter := bson.M{"_id": user.ID}
	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": user})
	return err
}

// UpdateUserProfile implements UserRepository.
func (u *UserMongoRepository) UpdateUserProfile(userID string, profile *UserProfile) error {
	ctx := context.Background()
	collection := u.client.Database(u.config.UsersDBName).Collection(UserProfileCollectionName)
	filter := bson.M{"_id": userID}
	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": profile})
	return err
}

// DeleteExpiredTokens implements UserRepository.
func (u *UserMongoRepository) DeleteExpiredTokens() error {
	ctx := context.Background()
	collection := u.client.Database(u.config.UsersDBName).Collection(UserTokenCollectionName)
	filter := bson.M{"expires_at": bson.M{"$lt": time.Now()}}
	_, err := collection.DeleteMany(ctx, filter)
	return err
}

// GetUserTokenByToken implements UserRepository.
func (u *UserMongoRepository) GetUserTokenByToken(token string) (*UserToken, error) {
	ctx := context.Background()
	collection := u.client.Database(u.config.UsersDBName).Collection(UserTokenCollectionName)
	filter := bson.M{"token": token, "is_valid": true}
	result := collection.FindOne(ctx, filter)

	if result.Err() != nil {
		return nil, result.Err()
	}

	var userToken UserToken
	err := result.Decode(&userToken)
	if err != nil {
		return nil, err
	}

	// Check if token is expired
	if time.Now().After(userToken.ExpiresAt) {
		return nil, errors.New("token is expired")
	}

	return &userToken, nil
}

// InvalidateToken implements UserRepository.
func (u *UserMongoRepository) InvalidateToken(userID string) error {
	ctx := context.Background()
	collection := u.client.Database(u.config.UsersDBName).Collection(UserTokenCollectionName)
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"is_valid": false}}
	_, err := collection.UpdateMany(ctx, filter, update)
	return err
}
