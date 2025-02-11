package model

import "go.mongodb.org/mongo-driver/mongo"

type CreateResult struct {
	Message string
	Result  *mongo.InsertOneResult
}
