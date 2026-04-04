package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hellomail/hellomail/internal/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	direction := flag.String("direction", "up", "migration direction: up or down")
	steps := flag.Int("steps", 0, "number of migrations to apply (0 = all)")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	fmt.Printf("Running migrations %s on %s\n", *direction, cfg.DatabaseURL[:40]+"...")

	// golang-migrate will be integrated here
	// For now, use: migrate -path internal/db/migrations -database $DATABASE_URL up
	_ = steps
	_ = direction

	fmt.Println("Note: Install golang-migrate CLI and run:")
	fmt.Printf("  migrate -path internal/db/migrations -database \"%s\" %s\n", cfg.DatabaseURL, *direction)

	return nil
}
