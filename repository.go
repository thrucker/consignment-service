package main

import (
	"context"
	pb "github.com/thrucker/consignment-service/proto/consignment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository interface {
	Create(consignment *pb.Consignment) error
	GetAll() ([]*pb.Consignment, error)
}

type MongoRepository struct {
	collection *mongo.Collection
}

func (repo *MongoRepository) Create(consignment *pb.Consignment) error {
	_, err := repo.collection.InsertOne(context.Background(), consignment)
	return err
}

func (repo *MongoRepository) GetAll() ([]*pb.Consignment, error) {
	cur, err := repo.collection.Find(context.Background(), bson.D{})

	if err != nil {
		return nil, err
	}

	var consignments []*pb.Consignment
	for cur.Next(context.Background()) {
		var consignment *pb.Consignment
		if err := cur.Decode(&consignment); err != nil {
			return nil, err
		}
		consignments = append(consignments, consignment)
	}
	return consignments, err
}
