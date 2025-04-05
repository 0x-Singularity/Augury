package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/microsoft/go-mssqldb" // Azure SQL driver
)

// QueryLog represents a lookup log entry
type QueryLog struct {
	LogID       int       `json:"log_id"`
	IOC         string    `json:"ioc"`
	LastLookup  time.Time `json:"last_lookup"`
	ResultCount int       `json:"result_count"`
	UserName    string    `json:"user_name"`
}

// DB connection (initialized once)
var db *sql.DB

// ConnectDB initializes the database connection
func ConnectDB() error {
	// Fetch connection parameters from environment variables
	server := os.Getenv("AZURE_SQL_SERVER")
	database := os.Getenv("AZURE_SQL_DATABASE")
	user := os.Getenv("AZURE_SQL_USER")
	password := os.Getenv("AZURE_SQL_PASSWORD")

	// Azure SQL connection string
	connStr := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&encrypt=true",
		user, password, server, database)

	// Open database connection
	var err error
	db, err = sql.Open("sqlserver", connStr)
	if err != nil {
		log.Println("Error opening database:", err)
		return err
	}

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Println("Error connecting to database:", err)
		return err
	}

	log.Println("Connected to Azure SQL successfully!")
	return nil
}

// InsertQueryLog adds a new IOC lookup record
func InsertQueryLog(ioc string, resultCount int, userName string) error {
	query := `INSERT INTO QueryLogs (ioc, last_lookup, result_count, user_name) 
	          VALUES (@ioc, GETDATE(), @resultCount, @userName)`

	_, err := db.Exec(query,
		sql.Named("ioc", ioc),
		sql.Named("resultCount", resultCount),
		sql.Named("userName", userName),
	)

	if err != nil {
		log.Println("Error inserting query log:", err)
	}
	return err
}

// GetQueryLog retrieves a log entry by IOC
func GetQueryLog(ioc string) ([]QueryLog, error) {
	query := `SELECT log_id, ioc, last_lookup, result_count, user_name
	FROM QueryLogs
	WHERE ioc = @ioc
	ORDER BY last_lookup DESC
	OFFSET 1 ROWS FETCH NEXT 1 ROWS ONLY;`

	rows, err := db.Query(query, sql.Named("ioc", ioc))
	if err != nil {
		log.Println("Error querying query logs:", err)
		return nil, err
	}
	defer rows.Close()

	var logs []QueryLog
	for rows.Next() {
		var logEntry QueryLog
		if err := rows.Scan(&logEntry.LogID, &logEntry.IOC, &logEntry.LastLookup, &logEntry.ResultCount, &logEntry.UserName); err != nil {
			log.Println("Error scanning query log row:", err)
			return nil, err
		}
		logs = append(logs, logEntry)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating query log rows:", err)
		return nil, err
	}

	// Return an empty slice instead of nil if no results found
	return logs, nil
}
