package main

import (
	"context"
	"fmt"
)

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
