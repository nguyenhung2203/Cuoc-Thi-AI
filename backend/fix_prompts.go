package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost port=15432 user=postgres password=postgres dbname=ai_interview sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	// Clear existing
	db.MustExec("DELETE FROM ai_request_logs;")
	db.MustExec("DELETE FROM ai_prompt_templates;")

	// Read and execute 16
	content16, err := os.ReadFile("migrations/000016_seed_ai_prompts.up.sql")
	if err != nil { log.Fatal(err) }
	db.MustExec(string(content16))

	// Read and execute 18
	content18, err := os.ReadFile("migrations/000018_seed_remaining_ai_prompts.up.sql")
	if err != nil { log.Fatal(err) }
	db.MustExec(string(content18))

	fmt.Println("Fixed prompts successfully!")
}
