package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/FlamestarRS/blogaggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {

		return fmt.Errorf("error: no username")
	}
	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Println("error: user does not exist")
		os.Exit(1)
	}
	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("Username has been set: " + username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {

		return fmt.Errorf("error: no username")
	}
	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		fmt.Println("error: user already exists")
		os.Exit(1)
	}
	params := database.CreateUserParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println("New user created:" + username)
	fmt.Println(user)
	return nil
}
