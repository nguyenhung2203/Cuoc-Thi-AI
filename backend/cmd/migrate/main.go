// Command migrate applies SQL migrations from backend/migrations.
//
// Usage:
//   go run cmd/migrate/main.go up      # apply all pending .up.sql
//   go run cmd/migrate/main.go down    # roll back the latest applied migration
//   go run cmd/migrate/main.go status  # show applied vs pending
//
// Migrations are plain files named NNN_name.up.sql / NNN_name.down.sql.
// Applied versions are tracked in the schema_migrations table.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"backend/internal/config"
)

var fileRe = regexp.MustCompile(`^(\d+)_(.+)\.(up|down)\.sql$`)

type migration struct {
	version  string
	name     string
	upPath   string
	downPath string
}

type migrationTracker struct {
	legacySingleRow bool
}

func main() {
	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	}
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	// Keep the application ledger separate from Supabase/migrate's legacy
	// schema_migrations(version bigint, dirty boolean) table.
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS public.app_schema_migrations (
		version TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`); err != nil {
		log.Fatalf("create app_schema_migrations: %v", err)
	}

	dir := migrationsDir()
	migs := loadMigrations(dir)

	switch cmd {
	case "up":
		runUp(db, migs)
	case "down":
		runDown(db, migs)
	case "status":
		showStatus(db, migs)
	default:
		log.Fatalf("unknown command %q (use up|down|status)", cmd)
	}
}

func migrationsDir() string {
	// Resolve relative to the working directory (backend/) or this file.
	candidates := []string{"migrations", filepath.Join("..", "migrations")}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	log.Fatal("migrations directory not found")
	return ""
}

func loadMigrations(dir string) []migration {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("read migrations: %v", err)
	}
	byKey := map[string]*migration{}
	for _, e := range entries {
		m := fileRe.FindStringSubmatch(e.Name())
		if m == nil {
			continue
		}
		version, name, dir2 := m[1], m[2], m[3]
		key := version + "_" + name
		mig := byKey[key]
		if mig == nil {
			mig = &migration{version: key, name: name}
			byKey[key] = mig
		}
		full := filepath.Join(dir, e.Name())
		if dir2 == "up" {
			mig.upPath = full
		} else {
			mig.downPath = full
		}
	}
	var out []migration
	for _, m := range byKey {
		out = append(out, *m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].version < out[j].version })
	return out
}

func applied(db *sqlx.DB) map[string]bool {
	set := map[string]bool{}

	var appRows []struct {
		Version string `db:"version"`
		Name    string `db:"name"`
	}
	if err := db.Select(&appRows, `SELECT version, name FROM public.app_schema_migrations`); err == nil {
		for _, r := range appRows {
			set[r.Version] = true
			set[r.Name] = true
			if parts := strings.SplitN(r.Version, "_", 2); len(parts) > 0 {
				set[parts[0]] = true
			}
		}
	}

	var legacyVersions []int64
	if err := db.Select(&legacyVersions, `SELECT version FROM public.schema_migrations WHERE dirty = false`); err == nil {
		for _, v := range legacyVersions {
			vStr := fmt.Sprintf("%06d", v)
			set[vStr] = true
		}
	}

	return set
}

func runUp(db *sqlx.DB, migs []migration) {
	if _, err := db.Exec("SELECT pg_advisory_lock(hashtext('cuocthi_ai_migrations'))"); err != nil {
		log.Printf("warning: acquire migration lock failed: %v", err)
	} else {
		defer db.Exec("SELECT pg_advisory_unlock(hashtext('cuocthi_ai_migrations'))")
	}

	done := applied(db)
	count := 0
	for _, m := range migs {
		prefix := strings.SplitN(m.version, "_", 2)[0]
		if done[m.version] || done[prefix] || done[m.name] {
			continue
		}
		if m.upPath == "" {
			log.Printf("skipping %s (no up file)", m.version)
			continue
		}
		sqlBytes, err := os.ReadFile(m.upPath)
		if err != nil {
			log.Fatalf("read %s: %v", m.upPath, err)
		}
		log.Printf("applying %s", m.version)

		rawSQL := string(sqlBytes)
		statements := splitSQLStatements(rawSQL)
		migFailed := false
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := db.Exec(stmt); err != nil {
				errMsg := strings.ToLower(err.Error())
				if strings.Contains(errMsg, "already exists") || strings.Contains(errMsg, "duplicate") || strings.Contains(errMsg, "multiple primary keys") {
					log.Printf("  [notice] %s: %v", m.version, err)
				} else {
					log.Printf("  [warning] %s stmt error: %v", m.version, err)
					migFailed = true
				}
			}
		}

		_, _ = db.Exec(`INSERT INTO public.app_schema_migrations (version, name) VALUES ($1, $2) ON CONFLICT DO NOTHING`, m.version, m.name)
		if !migFailed {
			count++
		}
	}
	log.Printf("up complete: %d migration(s) processed", count)
}

func splitSQLStatements(sql string) []string {
	var stmts []string
	for _, part := range strings.Split(sql, ";") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			stmts = append(stmts, trimmed)
		}
	}
	return stmts
}

func runDown(db *sqlx.DB, migs []migration) {
	if _, err := db.Exec("SELECT pg_advisory_lock(hashtext('cuocthi_ai_migrations'))"); err != nil {
		log.Fatalf("acquire migration lock: %v", err)
	}
	defer db.Exec("SELECT pg_advisory_unlock(hashtext('cuocthi_ai_migrations'))")

	done := applied(db)
	// Roll back the highest applied version only.
	for i := len(migs) - 1; i >= 0; i-- {
		m := migs[i]
		if !done[m.version] {
			continue
		}
		if m.downPath == "" {
			log.Fatalf("no down migration for %s", m.version)
		}
		sqlBytes, err := os.ReadFile(m.downPath)
		if err != nil {
			log.Fatalf("read %s: %v", m.downPath, err)
		}
		log.Printf("rolling back %s_%s", m.version, m.name)
		tx := db.MustBegin()
		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			log.Fatalf("rollback %s failed: %v", m.version, err)
		}
		if _, err := tx.Exec(`DELETE FROM public.app_schema_migrations WHERE version = $1`, m.version); err != nil {
			_ = tx.Rollback()
			log.Fatalf("remove migration record %s: %v", m.version, err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("commit rollback %s: %v", m.version, err)
		}
		log.Printf("down complete: rolled back %s", m.version)
		return
	}
	log.Println("down: nothing to roll back")
}

func showStatus(db *sqlx.DB, migs []migration) {
	done := applied(db)
	var b strings.Builder
	for _, m := range migs {
		status := "pending"
		if done[m.version] {
			status = "applied"
		}
		fmt.Fprintf(&b, "  [%s] %s_%s\n", status, m.version, m.name)
	}
	fmt.Printf("Migrations (%d):\n%s", len(migs), b.String())
}
