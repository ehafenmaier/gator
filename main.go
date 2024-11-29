package main

import (
	"github.com/ehafenmaier/boot-dev-gator/internal/config"
	"log"
	"os"
)

type state struct {
	config *config.Config
}

func main() {
	// Read the configuration
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new state
	s := &state{
		config: cfg,
	}

	// Create commands map
	c := commands{
		cmds: make(map[string]func(*state, command) error),
	}

	// Register the login command
	c.register("login", handlerLogin)

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
