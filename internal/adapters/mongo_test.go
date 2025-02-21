package adapters

import (
	"context"
	"testing"
	"time"

	"github.com/filipegorges/ports/internal/app/domain"
	"github.com/filipegorges/ports/internal/shared/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoRepository(t *testing.T) {
	uri := testutils.SetupMongoContainer(t)
	config := MongoConfig{
		URI:            uri,
		Database:       "testdb",
		Collection:     "ports",
		ConnectTimeout: 10 * time.Second,
	}

	tests := []struct {
		name            string
		repoCtx         context.Context
		upsertCtx       context.Context
		port            *domain.Port
		expectRepoErr   bool
		expectUpsertErr bool
	}{
		{
			name:      "Upsert Success",
			repoCtx:   context.Background(),
			upsertCtx: context.Background(),
			port: &domain.Port{
				Name:        "test port",
				City:        "test city",
				Country:     "test country",
				Alias:       []string{"alias1", "alias2"},
				Regions:     []string{"region"},
				Coordinates: []float64{12.34, 56.78},
				Province:    "test province",
				Timezone:    "UTC",
				Unlocs:      []string{"AAAAA"},
				Code:        "1234",
			},
			expectRepoErr:   false,
			expectUpsertErr: false,
		},
		{
			name:            "Upser with nil Port",
			repoCtx:         context.Background(),
			upsertCtx:       context.Background(),
			port:            nil,
			expectRepoErr:   false,
			expectUpsertErr: true,
		},
		{
			name:      "Upsert with cancelled context",
			repoCtx:   context.Background(),
			upsertCtx: func() context.Context { ctx, cancel := context.WithCancel(context.Background()); cancel(); return ctx }(),
			port: &domain.Port{
				Name:        "cancelled port",
				City:        "test city",
				Country:     "test country",
				Alias:       []string{"alias1"},
				Regions:     []string{"region"},
				Coordinates: []float64{98.76, 54.32},
				Province:    "test province",
				Timezone:    "UTC",
				Unlocs:      []string{"AAAAA"},
				Code:        "1234",
			},
			expectRepoErr:   false,
			expectUpsertErr: true,
		},
		{
			name: "Repository creation with cancelled context",
			repoCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			upsertCtx:       nil,
			port:            nil,
			expectRepoErr:   true,
			expectUpsertErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := NewMongoRepository(tt.repoCtx, config)
			if tt.expectRepoErr {
				require.Error(t, err)
				require.Nil(t, repo)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, repo)

			if tt.port != nil {
				err := repo.Upsert(tt.upsertCtx, tt.port)
				if tt.expectUpsertErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)

					var result domain.Port
					err = repo.collection.FindOne(
						context.Background(),
						bson.M{"coordinates": tt.port.Coordinates},
					).Decode(&result)
					require.NoError(t, err)
					assert.Equal(t, tt.port.Name, result.Name)
				}
			}
		})
	}
}
