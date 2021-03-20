package services

import (
	"github.com/alanwade2001/spa-submissions-api/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoIdGenerator struct{}

func NewMongoIdGenerator() types.IdGeneratorAPI {
	return &MongoIdGenerator{}
}

func (i *MongoIdGenerator) Next() string {
	return primitive.NewObjectID().Hex()
}
