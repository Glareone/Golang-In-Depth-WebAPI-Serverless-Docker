package database

// I use standard database/sql database API with Postgres PGX adapter
// it's less performant, but provides flexibility over real drivers you use behind the scenes
// another option is to use "github.com/jackc/pgx" API directly, it's more performant
import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib" // Import the adapter
	"log"
	"os"
)

// DB global variable of Database to be used elsewhere in this application
var DB *sql.DB

func InitDatabase() {
	var err error

	// Get environment variables
	var dbHost = os.Getenv("DB_HOST")
	var dbPort = os.Getenv("DB_PORT")
	var dbUser = os.Getenv("DB_USER")
	var dbPassword = os.Getenv("DB_PASSWORD")
	var dbName = os.Getenv("DB_NAME")

	// Construct connection string using environment variables
	var connectionStringPostgres = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	DB, err = sql.Open("pgx", connectionStringPostgres)
	if err != nil {
		// even if we panic here the server will not crash because of Gin
		panic("Connection with the database cannot be established")
	}

	// number of connections we keep opened while no one uses them
	DB.SetMaxIdleConns(5)

	// number of connections could simultaneously be opened
	DB.SetMaxOpenConns(10)

	// Create Tables
	migrationsCreateTables()

	// All Further modifications run here
	migrationsAlterTables()
}

func migrationsCreateTables() {
	var createTableUsers = `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			passwordHash TEXT NOT NULL
		)
	`

	var createTableEvents = `
		CREATE TABLE IF NOT EXISTS events (
        	id SERIAL PRIMARY KEY,
        	name TEXT NOT NULL,
            description TEXT NOT NULL,
            location TEXT NOT NULL,
            dateTime TIMESTAMP NOT NULL,
        	user_id INTEGER
		)
	`

	_, err := DB.Exec(createTableUsers)
	if err != nil {
		log.Fatalf("[Migration] Users table has not been created properly: %w", err)
	}

	_, err = DB.Exec(createTableEvents)
	if err != nil {
		log.Fatalf("[Migration] Events table has not been created properly: %w", err)
	}
}
