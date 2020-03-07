package consignmentMongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)


// MongoRepository implementation
type MongoRepository struct {
	Collection *mongo.Collection
}

func (repository *MongoRepository) Create(ctx context.Context, consignment *Consignment) error {
	_, err := repository.Collection.InsertOne(ctx, consignment)
	return err
}

func (repository *MongoRepository) GetAll(ctx context.Context) ([]*Consignment, error) {
	cur, err := repository.Collection.Find(ctx, nil, nil)
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

