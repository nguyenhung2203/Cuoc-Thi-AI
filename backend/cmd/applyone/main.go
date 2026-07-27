package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: applyone <sql-file> [--query]")
	}
	path := os.Args[1]
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read %s: %v", path, err)
	}
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer db.Close()
	if len(os.Args) > 2 && os.Args[2] == "--query" {
		rows, err := db.Queryx(string(b))
		if err != nil {
			log.Fatalf("query: %v", err)
		}
		defer rows.Close()
		for rows.Next() {
			m, _ := rows.SliceScan()
			fmt.Println(m...)
		}
		return
	}
	res, err := db.Exec(string(b))
	if err != nil {
		log.Fatalf("exec: %v", err)
	}
	n, _ := res.RowsAffected()
	fmt.Printf("OK (%d rows)\n", n)
}
