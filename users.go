package main

import (
	"context"
	"fmt"
	"github.com/ehafenmaier/boot-dev-gator/internal/database"
	"github.com/google/uuid"
	"time"
)

// Register user handler function
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

// Login user handler function
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: login <username>")
	}

	if len(cmd.args[0]) == 0 {
		return fmt.Errorf("username cannot be empty")
	}

	// Check if the user exists
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("user does not exist")
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Logged in as %s\n", cmd.args[0])

	return nil
}

// Reset users handler function
func handlerReset(s *state, _ command) error {
	// Reset the users table in the database
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting users in the database: %v", err)
	}

	fmt.Println("Database users reset")

	return nil
}

// List users handler function
func handlerUsers(s *state, _ command) error {
	// Get all users from the database
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users from the database: %v", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}
