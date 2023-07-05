package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/config"
)

const (
	dialect     = "pgx" // database driver dialect used for migration.
	fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable" // database connection string.
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError) // Command-line flag set.
	dir   = flags.String("dir", "migrations", "directory with migration files") // flag specifying directory with migration files.
)

func main() {
	// Get the current directory
	dirx, err := os.Getwd()
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Current directory:", dirx)

	// Set the usage function for the flags
	flags.Usage = usage

	// Parse the command-line arguments
	flags.Parse(os.Args[1:])

	// Get the remaining arguments after the flags
	args := flags.Args()

	// Check if no arguments were provided or the help flag is present
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		// Display the usage information
		flags.Usage()
		return
	}

	// Retrieve the command from the arguments
	command := args[0]

	// Load the configuration
	c := config.New()

	// Construct the database connection string
	dbString := fmt.Sprintf(fmtDBString, c.DB.RDBMS.Host, c.DB.RDBMS.Username, c.DB.RDBMS.Password, c.DB.RDBMS.DBName, c.DB.RDBMS.Port)

	// Open the database connection
	db, err2 := goose.OpenDBWithDriver(dialect, dbString)
	if err2 != nil {
		log.Fatalf(err2.Error())
	}

	defer func() {
		// Close the database connection
		if err := db.Close(); err != nil {
			log.Fatalf(err.Error())
		}
	}()

	// Run the goose command
	if err := goose.Run(command, db, *dir, args[1:]...); err != nil {
		log.Fatalf("migrate %v: %v", command, err)
	}
}

// usage displays the usage information for the migrate command.
func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate COMMAND
Examples:
    migrate status
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations`
)
