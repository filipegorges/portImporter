package ports

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockPortService struct {
	mock.Mock
}

func (m *mockPortService) Import(ctx context.Context, input io.Reader) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func TestCLI_Run(t *testing.T) {
	tests := []struct {
		name           string
		args           func(filePath string) []string
		setupFile      func() (string, func())
		importBehavior func(m *mockPortService)
		expectedErr    string
	}{
		{
			name: "Successful run",
			args: func(filePath string) []string {
				return []string{"", filePath}
			},
			setupFile: func() (string, func()) {
				tmp, err := os.CreateTemp("", "testfile*.json")
				require.NoError(t, err)
				tmp.Close()
				return tmp.Name(), func() { os.Remove(tmp.Name()) }
			},
			importBehavior: func(m *mockPortService) {
				m.On("Import",
					mock.Anything,
					mock.Anything,
				).Return(nil)
			},
			expectedErr: "",
		},
		{
			name: "Import returns error",
			args: func(filePath string) []string {
				return []string{"", filePath}
			},
			setupFile: func() (string, func()) {
				tmp, err := os.CreateTemp("", "testfile*.json")
				require.NoError(t, err)
				tmp.Close()
				return tmp.Name(), func() { os.Remove(tmp.Name()) }
			},
			importBehavior: func(m *mockPortService) {
				m.On("Import",
					mock.Anything,
					mock.Anything,
				).Return(errors.New("import error"))
			},
			expectedErr: "import error",
		},
		{
			name: "Insufficient arguments",
			args: func(_ string) []string {
				return []string{""}
			},
			setupFile: func() (string, func()) {
				return "", func() {}
			},
			importBehavior: func(m *mockPortService) {},
			expectedErr:    "insufficient arguments provided",
		},
		{
			name: "Failed to open file",
			args: func(_ string) []string {
				return []string{"", "nonexistent.json"}
			},
			setupFile: func() (string, func()) {
				return "nonexistent.json", func() {}
			},
			importBehavior: func(m *mockPortService) {},
			expectedErr:    "failed to open file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFile()
			defer cleanup()

			os.Args = tt.args(filePath)

			mockSvc := &mockPortService{}
			tt.importBehavior(mockSvc)

			cli := NewCLI(mockSvc)
			err := cli.Run(context.Background())
			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
			mockSvc.AssertExpectations(t)
		})
	}
}
