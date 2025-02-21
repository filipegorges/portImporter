package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/filipegorges/ports/internal/app/domain"
)

type jsonReader struct{}

func NewJsonReader() *jsonReader {
	return &jsonReader{}
}

// Stream using best-effort: only file-level errors (or context cancellation) are returned,
// any error decoding an individual port is logged and the port skipped, so as to not lose all
// data in case a few have issues.
func (f *jsonReader) Stream(ctx context.Context, in io.Reader, out chan<- *domain.Port) error {
	dec := json.NewDecoder(in)

	tok, err := dec.Token()
	if err != nil {
		return fmt.Errorf("failed to read opening token: %w", err)
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '{' {
		return fmt.Errorf("expected '{' as the opening token, got: %v", tok)
	}

	for dec.More() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		tok, err := dec.Token()
		if err != nil {
			return fmt.Errorf("failed to read key: %w", err)
		}
		key, ok := tok.(string)
		if !ok {
			log.Printf("expected key to be a string but got: %v", tok)
			continue
		}

		var port domain.Port
		if err := dec.Decode(&port); err != nil {
			log.Printf("failed to decode record for key %s: %v", key, err)
			continue
		}
		out <- &port
	}

	_, err = dec.Token()
	if err != nil {
		return fmt.Errorf("failed to read closing token: %w", err)
	}

	return nil
}
