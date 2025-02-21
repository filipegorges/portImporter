package ports

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
)

type portService interface {
	Import(ctx context.Context, input io.Reader) error
}

type cli struct {
	portService portService
}

func NewCLI(portService portService) *cli {
	return &cli{
		portService: portService,
	}
}

// TODO: take in file argument via flags
func (p *cli) Run(ctx context.Context) error {
	log.Println("running portImporter CLI")
	if len(os.Args) < 2 {
		return fmt.Errorf("insufficient arguments provided, please provide a valid JSON file")
	}

	filePath := os.Args[1]
	log.Printf("importing ports from %q\n", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer file.Close()

	err = p.portService.Import(ctx, file)
	if err != nil {
		log.Printf("failed to import: %v\n", err)
		return err
	}
	log.Println("import completed successfully")
	return nil
}
