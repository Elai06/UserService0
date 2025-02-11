package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const users = "UserService"
const service = "user"

//go:generate mockgen -destination=mocks/mock_user_cache.go -package=mocks userService/internal/repository IUserService
type IUserService interface {
	CreateUser(user Data) (*mongo.InsertOneResult, error)
	GetUserByID(id int64) (*Data, error)
	GetUsers() (*[]Data, error)
}

type Data struct {
	UserId int64  `json:"userId"`
	Name   string `json:"name"`
}

type Repository struct {
	client *mongo.Client
}

func NewService(url string) (*Repository, error) {
	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return &Repository{client: client}, err
}

func (ur *Repository) CreateUser(user Data) (*mongo.InsertOneResult, error) {
	userId, err := ur.getNextUserID()
	if err != nil {
		return nil, err
	}
	user.UserId = userId
	collection := ur.getCollection()
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("Inserted document with ID:", insertResult.InsertedID)

	return insertResult, nil
}

func (ur *Repository) GetUserByID(id int64) (*Data, error) {
	filter := map[string]interface{}{"userId": id}
	result := Data{}
	collection := ur.getCollection()
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Printf("Found document: %+v\n", result)

	return &result, nil
}

func (ur *Repository) GetUsers() (*[]Data, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result []Data
	collection := ur.getCollection()

	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &result); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &result, nil
}

func (ur *Repository) getNextUserID() (int64, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "userId", Value: -1}})

	collection := ur.getCollection()
	var lastUser Data
	err := collection.FindOne(context.TODO(), bson.D{}, opts).Decode(&lastUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 1, err
		}
		log.Fatal(err)
		return 0, err
	}

	return lastUser.UserId + 1, nil
}

func (ur *Repository) getCollection() *mongo.Collection {
	return ur.client.Database(users).Collection(service)
}
