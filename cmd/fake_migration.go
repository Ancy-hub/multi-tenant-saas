package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	db, err := pgxpool.New(context.Background(), "postgres://postgres:root@localhost:8080/postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	// Create the schema_migrations table manually
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version bigint NOT NULL PRIMARY KEY,
			dirty boolean NOT NULL
		)
	`)
	if err != nil {
		log.Fatal("Failed to create schema_migrations:", err)
	}

	// Insert the baseline version 1
	_, err = db.Exec(ctx, `
		INSERT INTO schema_migrations (version, dirty)
		VALUES (1, false)
		ON CONFLICT (version) DO NOTHING
	`)
	if err != nil {
		log.Fatal("Failed to insert fake migration:", err)
	}

	fmt.Println("Successfully inserted fake migration record for version 1!")
}
