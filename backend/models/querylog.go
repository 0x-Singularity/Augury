package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// QueryLog represents a lookup log entry
type QueryLog struct {
	LogID       int       `json:"log_id"`
	IOC         string    `json:"ioc"`
	LastLookup  time.Time `json:"last_lookup"`
	ResultCount int       `json:"result_count"`
	UserName    string    `json:"user_name"`
}

var db *sql.DB

// ConnectDB initializes the PostgreSQL connection
func ConnectDB() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	//pool tuning
	db.SetMaxOpenConns(20)                  // total connections allowed (0 = unlimited)
	db.SetMaxIdleConns(10)                  // idle (kept‑alive) connections
	db.SetConnMaxIdleTime(30 * time.Minute) // recycle idle conns after 30 min
	db.SetConnMaxLifetime(2 * time.Hour)
	if err = db.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully!")
	return nil
}

// InsertQueryLog adds a new IOC lookup record
func InsertQueryLog(ioc string, resultCount int, userName string) error {
	const stmt = `
		INSERT INTO ioc_query_log (ioc, last_lookup, result_count, user_name)
		VALUES ($1, now(), $2, $3);
	`
	_, err := db.Exec(stmt, ioc, resultCount, userName)
	if err != nil {
		return fmt.Errorf("insert query log: %w", err)
	}
	return nil
}

// GetQueryLog returns the most‑recent previous log for an IOC
func GetQueryLog(ioc string) ([]QueryLog, error) {
	const stmt = `
		SELECT id, ioc, last_lookup, result_count, user_name
		FROM   ioc_query_log
		WHERE  ioc = $1
		ORDER  BY last_lookup DESC
		OFFSET 1           -- skip latest
		LIMIT  1;          -- return the previous
	`

	rows, err := db.Query(stmt, ioc)
	if err != nil {
		return nil, fmt.Errorf("select logs: %w", err)
	}
	defer rows.Close()

	var logs []QueryLog
	for rows.Next() {
		var q QueryLog
		if err := rows.Scan(&q.LogID, &q.IOC, &q.LastLookup, &q.ResultCount, &q.UserName); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		logs = append(logs, q)
	}
	return logs, rows.Err()
}
