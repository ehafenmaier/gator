package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, _ command) error {
	// Reset the users table in the database
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting users in the database: %v", err)
	}

	fmt.Println("Database users reset")

	return nil
}
