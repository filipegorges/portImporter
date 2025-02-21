package adapters

import (
	"bytes"
	"context"
	"testing"

	"github.com/filipegorges/ports/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonReader_Stream(t *testing.T) {
	validJSON := `{
		"port1": {"city": "1", "name": "port1"},
		"port2": {"city": "2", "name": "port2"}
	}`

	tests := []struct {
		name                 string
		input                string
		cancelContext        bool // if true, the context will be canceled immediately
		expectedPorts        []*domain.Port
		expectedErrSubstring string
	}{
		{
			name:          "Valid JSON",
			input:         validJSON,
			cancelContext: false,
			expectedPorts: []*domain.Port{
				{City: "1", Name: "port1"},
				{City: "2", Name: "port2"},
			},
		},
		{
			name:                 "Invalid opening token",
			input:                `[]`,
			cancelContext:        false,
			expectedErrSubstring: "expected '{' as the opening token",
			expectedPorts:        []*domain.Port(nil),
		},
		{
			name:          "Port decoding error",
			input:         `{"port1": "invalid"}`,
			cancelContext: false,
			expectedPorts: []*domain.Port(nil),
		},
		{
			name:                 "Missing closing token",
			input:                `{"port1": {"city": "1", "name": "port1"}`,
			cancelContext:        false,
			expectedErrSubstring: "failed to read closing token",
			expectedPorts: []*domain.Port{
				{City: "1", Name: "port1"},
			},
		},
		{
			name:                 "Empty input",
			input:                ``,
			cancelContext:        false,
			expectedErrSubstring: "failed to read opening token",
			expectedPorts:        []*domain.Port(nil),
		},
		{
			name:                 "Context cancelled",
			input:                validJSON,
			cancelContext:        true,
			expectedErrSubstring: "context canceled",
			expectedPorts:        []*domain.Port(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx context.Context
			if tt.cancelContext {
				c, cancel := context.WithCancel(context.Background())
				cancel()
				ctx = c
			} else {
				ctx = context.Background()
			}

			in := bytes.NewBufferString(tt.input)
			out := make(chan *domain.Port)
			errCh := make(chan error, 1)

			reader := NewJsonReader()

			go func() {
				defer close(errCh)
				defer close(out)
				errCh <- reader.Stream(ctx, in, out)
			}()

			var ports []*domain.Port
			for p := range out {
				ports = append(ports, p)
			}
			// TODO: sort ports for better stability
			assert.Equal(t, tt.expectedPorts, ports)

			err := <-errCh
			if tt.expectedErrSubstring != "" {
				require.ErrorContains(t, err, tt.expectedErrSubstring)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
