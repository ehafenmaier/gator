package main

import (
	"context"
	"fmt"
	"github.com/ehafenmaier/boot-dev-gator/internal/database"
	"github.com/google/uuid"
	"time"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	// Get the current user from the database
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}

	// Create feed in the database
	dbParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), dbParams)
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}

	fmt.Printf("Feed Added\nName: %s\nUrl: %s\n", feed.Name, feed.Url)

	return nil
}
