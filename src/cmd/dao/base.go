package dao

import "go.mongodb.org/mongo-driver/mongo"

type baseDao struct {
	client *mongo.Client
}

const DATABASE_NAME = "ms-payment"
