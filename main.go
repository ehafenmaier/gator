package main

import (
	"database/sql"
	"github.com/ehafenmaier/boot-dev-gator/internal/config"
	"github.com/ehafenmaier/boot-dev-gator/internal/database"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	// Read the configuration
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading configuration: %v", err)
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("error opening database connection: %v", err)
	}

	// Create a new application state
	s := &state{
		db:  database.New(db),
		cfg: cfg,
	}

	// Create commands map
	c := commands{
		cmds: make(map[string]func(*state, command) error),
	}

	// Register commands
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)

	// Check for the proper number of arguments
	if len(os.Args) < 2 {
		log.Fatal("usage: boot-dev-gator <command> [args...]")
	}

	// Split the command line arguments into a command and its arguments
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	// Run the command
	err = c.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
