package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Migrate database
	Migrate() error

	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	BulkInsertLogs(ctx context.Context, logs []Log) error
	GetProjectLogs(ctx context.Context, projectID string, limit int) ([]Log, error)
}

type Log struct {
	LogID     string
	ProjectID string
	Log       string
	Timestamp int64
}

type service struct {
	db *sql.DB
}

var (
	aetherMainDb = os.Getenv("MAIN_DB_DATABASE")
	database     = os.Getenv("DB_DATABASE")
	password     = os.Getenv("DB_PASSWORD")
	username     = os.Getenv("DB_USERNAME")
	port         = os.Getenv("DB_PORT")
	host         = os.Getenv("DB_HOST")
	dbInstance   *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := connectToDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func connectToDatabase() (*sql.DB, error) {
	// Attempt to connect to the aether-logs database
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Check if the connection to the aether-logs database is successful
	err = db.Ping()
	if err == nil {
		return db, nil
	}

	// If the aether-logs database does not exist, connect to the aether database and create it
	db.Close()
	connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, aetherMainDb)
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to aether database: %w", err)
	}

	_, err = db.Exec(`CREATE DATABASE "aether-logs"`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create aether-logs database: %w", err)
	}

	// Close the connection to the aether database
	db.Close()

	// Reconnect to the aether-logs database
	connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to aether-logs database: %w", err)
	}

	return db, nil
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

func (s *service) Migrate() error {
	migrationPath := "file://./migrations"

	driver, err := postgres.WithInstance(s.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

func (s *service) BulkInsertLogs(ctx context.Context, logs []Log) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO logs (projectId, log, timestamp)
        VALUES ($1, $2, $3)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, log := range logs {
		_, err = stmt.ExecContext(ctx, log.ProjectID, log.Log, log.Timestamp)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *service) GetProjectLogs(ctx context.Context, projectID string, limit int) ([]Log, error) {
	rows, err := s.db.QueryContext(ctx, `
        SELECT logid, projectId, log, timestamp
        FROM logs
        WHERE projectId = $1
        ORDER BY timestamp DESC
        LIMIT $2
    `, projectID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []Log
	for rows.Next() {
		var log Log
		err := rows.Scan(&log.LogID, &log.ProjectID, &log.Log, &log.Timestamp)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}
