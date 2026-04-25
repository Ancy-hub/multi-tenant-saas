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

	tables := []string{"users", "organizations", "memberships", "projects", "tasks"}

	for _, table := range tables {
		fmt.Printf("--- TABLE: %s ---\n", table)
		rows, err := db.Query(ctx, "SELECT column_name, data_type, character_maximum_length, is_nullable, column_default FROM information_schema.columns WHERE table_name = $1 ORDER BY ordinal_position", table)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var colName, dataType, isNullable string
			var maxLen, colDefault *string
			
			err := rows.Scan(&colName, &dataType, &maxLen, &isNullable, &colDefault)
			if err != nil {
				log.Fatal(err)
			}
			
			lenStr := ""
			if maxLen != nil {
				lenStr = fmt.Sprintf("(%s)", *maxLen)
			}
			
			defStr := ""
			if colDefault != nil {
				defStr = fmt.Sprintf(" DEFAULT %s", *colDefault)
			}
			
			nullStr := ""
			if isNullable == "NO" {
				nullStr = " NOT NULL"
			}
			
			fmt.Printf("%s %s%s%s%s\n", colName, dataType, lenStr, defStr, nullStr)
		}
		rows.Close()
		fmt.Println()
	}
}
