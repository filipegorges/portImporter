package app_test

import (
	"context"
	"os"
	"testing"

	"github.com/filipegorges/ports/internal/app"
	"github.com/filipegorges/ports/internal/shared/testutils"
	"github.com/stretchr/testify/require"
)

func Test_Application(t *testing.T) {
	ctx := context.Background()
	os.Args = []string{"", "../../resources/ports.json"}
	uri := testutils.SetupMongoContainer(t)
	os.Setenv("DATABASE_URI", uri)
	os.Setenv("DATABASE_CONNECTION_TIMEOUT_IN_SECONDS", "10")

	err := app.Run(ctx)
	require.NoError(t, err)
}
