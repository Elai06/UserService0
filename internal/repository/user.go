package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const users = "UserService"
const service = "user"

//go:generate mockgen -destination=mocks/mock_user_cache.go -package=mocks userService/internal/repository UserService
type UserService interface {
	CreateUser(ctx context.Context, user Data) (*mongo.InsertOneResult, error)
	GetUserByID(ctx context.Context, id int64) (*Data, error)
	GetUsers(ctx context.Context) (*[]Data, error)
}

type Data struct {
	UserID int64  `json:"userId"`
	Name   string `json:"name"`
}

type Repository struct {
	client *mongo.Client
}

var writeTimeOut time.Duration

func NewService(url string, writeTimeout time.Duration) (*Repository, error) {
	clientOptions := options.Client().ApplyURI(url)
	writeTimeOut = writeTimeout

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("connect mongodb error: %w", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("Connected to MongoDB!")

	return &Repository{client: client}, err
}

func (ur *Repository) CreateUser(ctx context.Context, user Data) (*mongo.InsertOneResult, error) {
	userID, err := ur.getNextUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get next user id: %w", err)
	}

	user.UserID = userID
	collection := ur.getCollection()

	insertResult, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user into MongoDB: %w", err)
	}

	fmt.Println("Inserted document with ID:", insertResult.InsertedID)

	return insertResult, nil
}

func (ur *Repository) GetUserByID(ctx context.Context, id int64) (*Data, error) {
	filter := map[string]interface{}{"userId": id}
	result := Data{}
	collection := ur.getCollection()

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	log.Printf("Found document: %+v\n", result)

	return &result, nil
}

func (ur *Repository) GetUsers(ctx context.Context) (*[]Data, error) {
	ctx, cancel := context.WithTimeout(ctx, writeTimeOut*time.Second)
	defer cancel()

	var result []Data

	collection := ur.getCollection()

	cursor, errCollection := collection.Find(ctx, bson.M{})
	if errCollection != nil {
		return nil, fmt.Errorf("failed to find users: %w", errCollection)
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &result); err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	return &result, nil
}

func (ur *Repository) getNextUserID(ctx context.Context) (int64, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "userId", Value: -1}})
	collection := ur.getCollection()

	var lastUser Data
	err := collection.FindOne(ctx, bson.D{}, opts).Decode(&lastUser)
	if err != nil {
		return 0, fmt.Errorf("failed to find next user id: %w", err)
	}

	return lastUser.UserID + 1, nil
}

func (ur *Repository) getCollection() *mongo.Collection {
	return ur.client.Database(users).Collection(service)
}
