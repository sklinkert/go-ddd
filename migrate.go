package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	var databaseURL = flag.String("database-url", "", "Database connection URL")
	var migrationsPath = flag.String("migrations-path", "file://migrations", "Path to migrations directory")
	var command = flag.String("command", "up", "Migration command: up, down, version, force")
	var steps = flag.Int("steps", -1, "Number of migration steps (for up/down commands)")
	var version = flag.Int("version", -1, "Target version (for force command)")

	flag.Parse()

	if *databaseURL == "" {
		*databaseURL = os.Getenv("DATABASE_URL")
		if *databaseURL == "" {
			log.Fatal("Database URL is required. Use -database-url flag or set DATABASE_URL environment variable")
		}
	}

	db, err := sql.Open("pgx", *databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create database driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(*migrationsPath, "postgres", driver)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}
	defer m.Close()

	switch *command {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration up failed:", err)
		}
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to apply")
		} else {
			fmt.Println("Migrations applied successfully")
		}

	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration down failed:", err)
		}
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to rollback")
		} else {
			fmt.Println("Migrations rolled back successfully")
		}

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal("Failed to get version:", err)
		}
		fmt.Printf("Current version: %d (dirty: %v)\n", version, dirty)

	case "force":
		if *version < 0 {
			log.Fatal("Version is required for force command. Use -version flag")
		}
		err = m.Force(*version)
		if err != nil {
			log.Fatal("Force migration failed:", err)
		}
		fmt.Printf("Forced migration to version %d\n", *version)

	default:
		log.Fatal("Unknown command:", *command)
	}
}
