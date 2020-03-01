package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)


// MongoRepository implementation
type MongoRepository struct {
	collection *mongo.Collection
}

func (repository *MongoRepository) Create(ctx context.Context, consignment *Consignment) error {
	_, err := repository.collection.InsertOne(ctx, consignment)
	return err
}

func (repository *MongoRepository) GetAll(ctx context.Context) ([]*Consignment, error) {
	cur, err := repository.collection.Find(ctx, nil, nil)
	var consignments []*Consignment
	for cur.Next(ctx) {
		var consignment *Consignment
		if err := cur.Decode(&consignment); err != nil {
			return nil, err
		}
		consignments = append(consignments, consignment)
	}
	return consignments, err
}

