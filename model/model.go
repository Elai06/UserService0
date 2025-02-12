package model

import "go.mongodb.org/mongo-driver/mongo"

type ResponseUserService struct {
	Message string
	Result  *mongo.InsertOneResult
}
