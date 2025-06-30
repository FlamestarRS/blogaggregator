package main

import (
	"fmt"
	"os"

	"github.com/FlamestarRS/blogaggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	s := state{
		cfg: &cfg,
	}
	commandsMap := make(map[string]func(*state, command) error)
	commands := commands{
		cmds: commandsMap,
	}

	commands.register("login", handlerLogin)

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

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
}
