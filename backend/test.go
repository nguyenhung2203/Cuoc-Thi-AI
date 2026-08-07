package main

import (
	"context"
	"fmt"
	"log"

	"backend/internal/repository"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "host=localhost port=15432 user=postgres password=postgres dbname=ai_interview sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()


	repo2 := repository.NewCandidatePortalRepository(db)
	userID := "40408bac-a906-4684-b7cb-ec90997d1938"
	apps, _ := repo2.GetApplicationsByUserID(context.Background(), userID)
	fmt.Printf("Apps: %+v\n", apps)
}
