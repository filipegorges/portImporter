package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/filipegorges/ports/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPortRepository struct {
	mock.Mock
}

func (m *mockPortRepository) Upsert(ctx context.Context, port *domain.Port) error {
	args := m.Called(ctx, port)
	return args.Error(0)
}

type mockReader struct {
	mock.Mock
}

func (m *mockReader) Stream(ctx context.Context, file io.Reader, out chan<- *domain.Port) error {
	call := m.Called(ctx, file, out)
	return call.Error(0)
}

func TestImport(t *testing.T) {
	tests := []struct {
		name                string
		streamError         error
		ports               []*domain.Port
		upsertReturnErrors  []error
		expectedErrorSubstr string
	}{
		{
			name:        "successful import",
			streamError: nil,
			ports: []*domain.Port{
				{Name: "port1"},
				{Name: "port2"},
			},
			upsertReturnErrors:  []error{nil, nil},
			expectedErrorSubstr: "",
		},
		{
			name:        "successful import with upsert error",
			streamError: nil,
			ports: []*domain.Port{
				{Name: "port1"},
				{Name: "port2"},
			},
			upsertReturnErrors:  []error{nil, fmt.Errorf("upsert failure")},
			expectedErrorSubstr: "",
		},
		{
			name:                "stream error with no ports",
			streamError:         fmt.Errorf("stream failure"),
			ports:               nil,
			upsertReturnErrors:  nil,
			expectedErrorSubstr: "error streaming ports: stream failure",
		},
		{
			name:        "stream error after sending ports",
			streamError: fmt.Errorf("stream error"),
			ports: []*domain.Port{
				{Name: "port1"},
				{Name: "port2"},
			},
			upsertReturnErrors:  []error{nil, nil},
			expectedErrorSubstr: "error streaming ports: stream error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			input := strings.NewReader("dummy")
			repo := new(mockPortRepository)
			reader := new(mockReader)
			reader.
				On("Stream", ctx, input, mock.Anything).
				Return(tt.streamError).
				Run(func(args mock.Arguments) {
					const channelParameter = 2
					out := args.Get(channelParameter).(chan<- *domain.Port)
					if tt.ports != nil {
						for _, p := range tt.ports {
							out <- p
						}
					}
				})
			if tt.ports != nil {
				for i, port := range tt.ports {
					repo.On("Upsert", ctx, port).Return(tt.upsertReturnErrors[i])
				}
			}
			importer := NewportImporter(repo, reader)
			err := importer.Import(ctx, input)
			if tt.expectedErrorSubstr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorSubstr)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
			reader.AssertExpectations(t)
		})
	}
}
