package main

import (
	"context"
	"fmt"
	"github.com/ehafenmaier/boot-dev-gator/internal/database"
	"github.com/google/uuid"
	"time"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: register <username>")
	}

	if len(cmd.args[0]) == 0 {
		return fmt.Errorf("username cannot be empty")
	}

	// Check if the user already exists
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err == nil {
		return fmt.Errorf("user already exists")
	}

	// Create user in the database
	dbParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	user, err := s.db.CreateUser(context.Background(), dbParams)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	// Set current user in the configuration
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("Registered as %s\n", user.Name)

	return nil
}
