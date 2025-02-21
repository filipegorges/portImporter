package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/filipegorges/ports/internal/app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI            string
	Database       string
	Collection     string
	ConnectTimeout time.Duration
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(ctx context.Context, config MongoConfig) (*MongoRepository, error) {
	client, err := connectToMongo(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo repository: %w", err)
	}

	return &MongoRepository{
		collection: client.Database(config.Database).Collection(config.Collection),
	}, nil
}

func (r *MongoRepository) Upsert(ctx context.Context, port *domain.Port) error {
	if port == nil {
		return fmt.Errorf("port cannot be nil")
	}

	// using coordinates as unlocs being an array threw me off on its potential stability
	filter := bson.M{"coordinates": port.Coordinates}
	update := bson.M{"$set": port}
	opts := options.Update().SetUpsert(true)

	// TODO: add retrying logic
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert port: %w", err)
	}

	return nil
}

func connectToMongo(ctx context.Context, config MongoConfig) (*mongo.Client, error) {
	mongoCtx, cancel := context.WithTimeout(ctx, config.ConnectTimeout)
	defer cancel()

	clientOpts := options.Client().ApplyURI(config.URI)
	client, err := mongo.Connect(mongoCtx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	if err := client.Ping(mongoCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	return client, nil
}
