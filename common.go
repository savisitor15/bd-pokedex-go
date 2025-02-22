package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]cliCommand{
return map[string]cliCommand{
	"help" : {
		name: "help",
		description: "Displays a help message",
		callback: commandHelp,
	},
	"exit" : {
		name: "exit",
		description: "Exits the pokedex",
		callback: commandExit,
	},
}
}

func commandLoop() error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() != true {
			return scanner.Err()
		}
		input := cleanInput(scanner.Text())
		fmt.Println()
		if cmd, ok := getCommands()[input[0]]; !ok {
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

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, elm := range getCommands() {
		fmt.Printf("%s : %s\n", elm.name, elm.description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
