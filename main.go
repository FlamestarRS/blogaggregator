package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/FlamestarRS/blogaggregator/internal/config"
	"github.com/FlamestarRS/blogaggregator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	s := state{
		cfg: &cfg,
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	s.db = dbQueries

	commandsMap := make(map[string]func(*state, command) error)
	commands := commands{
		cmds: commandsMap,
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerListUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerListFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))

	input := os.Args
	if len(input) < 2 {
		fmt.Println("Error: No command")
		os.Exit(1)
	}

	cmdName := input[1]
	cmdArgs := []string{}
	if len(input) > 2 {
		cmdArgs = input[2:]
	}
	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}
	err = commands.run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
