package service

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/filipegorges/ports/internal/app/domain"
)

type portRepository interface {
	Upsert(ctx context.Context, port *domain.Port) error
}

type reader interface {
	Stream(ctx context.Context, file io.Reader, out chan<- *domain.Port) error
}

type portImporter struct {
	repository portRepository
	reader     reader
}

func NewportImporter(repo portRepository, r reader) *portImporter {
	return &portImporter{
		repository: repo,
		reader:     r,
	}
}

func (s *portImporter) Import(ctx context.Context, input io.Reader) error {
	ports := make(chan *domain.Port)
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)
		defer close(ports)

		errCh <- s.reader.Stream(ctx, input, ports)
	}()

	// TODO: batch to lower IO
	for port := range ports {
		if err := s.repository.Upsert(ctx, port); err != nil {
			log.Printf("failed to upsert port %q: %v", port.Name, err)
		}
	}

	if err := <-errCh; err != nil {
		return fmt.Errorf("error streaming ports: %w", err)
	}
	return nil
}
