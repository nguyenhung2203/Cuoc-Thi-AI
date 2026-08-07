package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
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

	companyID := "aaaaaaaa-0000-0000-0000-000000000001"

	// Fetch valid template ID if any, or NULL
	var templateID sql.NullString
	err = db.QueryRow("SELECT id FROM ai_prompt_templates LIMIT 1").Scan(&templateID)
	if err != nil {
		log.Printf("No template found, using NULL: %v", err)
	}

	actions := []struct {
		Name      string
		TokensIn  int
		TokensOut int
		Latency   int
	}{
		{"analyze_cv", 1450, 680, 1850},
		{"analyze_jd", 1200, 520, 1420},
		{"match_candidate", 980, 410, 1100},
		{"generate_questions", 850, 750, 2100},
		{"evaluate_interview", 3200, 1400, 3800},
		{"mock_interview_turn", 620, 290, 890},
	}

	stmt, err := db.Prepare(`
		INSERT INTO ai_request_logs (
			company_id, template_id, template_version,
			input_json, output_json, latency_ms,
			tokens_in, tokens_out, cost, status, created_at
		) VALUES (
			$1::uuid, $2, 1,
			$3::jsonb, '{"status":"success"}'::jsonb, $4,
			$5, $6, $7, 'success', NOW() - ($8 * INTERVAL '1 day')
		)
	`)
	if err != nil {
		log.Fatalf("Failed to prepare log insert: %v", err)
	}
	defer stmt.Close()

	totalLogs := 0
	var totalIn, totalOut int64

	for day := 30; day >= 0; day-- {
		numReqs := rand.Intn(5) + 4
		for r := 0; r < numReqs; r++ {
			act := actions[rand.Intn(len(actions))]
			tin := act.TokensIn + rand.Intn(400) - 200
			tout := act.TokensOut + rand.Intn(200) - 100
			lat := act.Latency + rand.Intn(400) - 200
			cost := float64(tin+tout) * 0.00000015
			inputJSON := fmt.Sprintf(`{"action":"%s","tokens":%d}`, act.Name, tin+tout)

			var tID interface{}
			if templateID.Valid {
				tID = templateID.String
			} else {
				tID = nil
			}

			_, err := stmt.Exec(companyID, tID, inputJSON, lat, tin, tout, cost, day)
			if err != nil {
				log.Printf("Insert error: %v", err)
			} else {
				totalLogs++
				totalIn += int64(tin)
				totalOut += int64(tout)
			}
		}
	}

	fmt.Printf("Successfully inserted %d AI request log records into Supabase PostgreSQL DB!\nTotal Tokens: %d (In: %d, Out: %d)\n",
		totalLogs, totalIn+totalOut, totalIn, totalOut)
}
