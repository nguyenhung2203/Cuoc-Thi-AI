package main

import (
	"context"
	"fmt"
	"log"

	"backend/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:15432/ai_interview?sslmode=disable")
	if err != nil {
		log.Fatalln("Connect Error:", err)
	}
	defer db.Close()

	var room models.InterviewRoom
	err = db.GetContext(context.Background(), &room, "SELECT * FROM interview_rooms WHERE interview_id = '10aee5ce-98fc-4dd0-8d1e-6b9dda50bcb0'")
	if err != nil {
		log.Fatalln("Get Error:", err)
	}
	fmt.Printf("RoomCode: %s\n", room.RoomCode)
}
