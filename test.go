package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type InterviewRoom struct {
	ID               string  `db:"id"`
	InterviewID      string  `db:"interview_id"`
	RoomCode         string  `db:"room_code"`
	Status           string  `db:"status"`
	Provider         *string `db:"provider"`
	ConnectionConfig []byte  `db:"connection_config"`
	OpenedAt         *string `db:"opened_at"`
	ClosedAt         *string `db:"closed_at"`
	CreatedAt        string  `db:"created_at"`
	UpdatedAt        string  `db:"updated_at"`
}

func main() {
	db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:5432/ai_interview?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	var room InterviewRoom
	err = db.GetContext(context.Background(), &room, "SELECT * FROM interview_rooms WHERE interview_id = '10aee5ce-98fc-4dd0-8d1e-6b9dda50bcb0'")
	if err != nil {
		log.Fatalln("Error:", err)
	}
	fmt.Printf("RoomCode: %s\n", room.RoomCode)
}
