package user

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Data struct {
	UserId int64  `json:"userId"`
	Name   string `json:"name"`
}

var collection *mongo.Collection

func ConnectToMongo() {
	url := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	collection = client.Database("UserService").Collection("user")
}

func CreateUser(user Data) *mongo.InsertOneResult {
	user.UserId = getNextUserID()
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted document with ID:", insertResult.InsertedID)

	return insertResult
}

func GetUserByID(id int64) Data {
	filter := map[string]interface{}{"userId": id}

	result := Data{}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return Data{}
	}

	fmt.Printf("Found document: %+v\n", result)

	return result
}

func GetUsers() []Data {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result []Data

	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}

	return result
}

func getNextUserID() int64 {
	opts := options.FindOne().SetSort(bson.D{{Key: "userId", Value: -1}})

	var lastUser Data
	err := collection.FindOne(context.TODO(), bson.D{}, opts).Decode(&lastUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1
		}
		log.Fatal(err)
	}

	return lastUser.UserId + 1
}
