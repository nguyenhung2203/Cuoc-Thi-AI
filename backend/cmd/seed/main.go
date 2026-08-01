package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgres.ouqlvqitktzqggwbmfki:khoihunglai@aws-1-ap-northeast-2.pooler.supabase.com:5432/postgres"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	sqlBytes, err := os.ReadFile("seeds/seed_demo.sql")
	if err != nil {
		log.Fatalf("Failed to read seed file: %v", err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("Failed to execute seed SQL: %v", err)
	}

	fmt.Println("Successfully executed seed_demo.sql into Supabase PostgreSQL DB!")
}
