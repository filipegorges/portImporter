package app

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/filipegorges/ports/internal/adapters"
	"github.com/filipegorges/ports/internal/app/service"
	"github.com/filipegorges/ports/internal/ports"
	"github.com/filipegorges/ports/internal/shared/env"
)

const (
	dbName         = "portImporter"
	collectionName = "ports"
)

func Run(ctx context.Context) error {
	log.Println("starting portImporter application")

	// TODO: create an .env file and load it through godotenv
	dbURI := env.MustEnv("DATABASE_URI")
	dbConnTimeout := env.MustEnv("DATABASE_CONNECTION_TIMEOUT_IN_SECONDS")
	connTimeout, err := strconv.Atoi(dbConnTimeout)
	if err != nil {
		log.Printf("database timeout must be an integer, found %q", dbConnTimeout)
		return err
	}

	repo, err := adapters.NewMongoRepository(ctx, adapters.MongoConfig{
		URI:            dbURI,
		Database:       dbName,
		Collection:     collectionName,
		ConnectTimeout: time.Second * time.Duration(connTimeout),
	})
	if err != nil {
		log.Printf("failed to initialize mongo repository: %v", err)
		return err
	}
	reader := adapters.NewJsonReader()
	svc := service.NewportImporter(repo, reader)
	cli := ports.NewCLI(svc)
	return cli.Run(ctx)
}
