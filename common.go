package main

import (
	"strings"
	"fmt"
	"bufio"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}



func commandLoop() error {
	commands := map[string]cliCommand{
		"exit" : {
			name: "exit",
			description: "Exits the pokedex",
			callback: commandExit,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() != true {
			return scanner.Err()
		}
		input := cleanInput(scanner.Text())
		fmt.Println()
		if cmd, ok := commands[input[0]]; !ok {
			fmt.Println("Unknown command!")
			fmt.Println("Your command was:", input[0])
			continue
		}else{
			cmd.callback()
		}
	}
}

func cleanInput(text string) []string {
	splits := strings.Split(text, " ")
	if len(splits) == 0 {
		return []string{}
	}
	final := make([]string, 0, len(splits))
	for _, spl := range splits {
		if len(strings.TrimSpace(spl)) == 0 {
			continue
		}
		final = append(final, strings.ToLower(strings.TrimSpace(spl)))
	}
	return final
}

func commandExit() error {
	os.Exit(0)
	return nil
}
