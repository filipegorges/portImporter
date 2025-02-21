package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupMongoContainer(t *testing.T) string {
	t.Helper()
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "latest", // TODO: pin down the same version of mongo used in production
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		err = pool.Purge(resource)
		if err != nil {
			t.Fatalf("Could not purge resource: %s", err)
		}
	})

	port := resource.GetPort("27017/tcp")
	uri := fmt.Sprintf("mongodb://localhost:%s", port)

	err = pool.Retry(func() error {
		clientOpts := options.Client().ApplyURI(uri)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, clientOpts)
		if err != nil {
			return err
		}
		return client.Ping(ctx, nil)
	})
	require.NoError(t, err)

	return uri
}
