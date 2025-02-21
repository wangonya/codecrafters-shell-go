package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

var validCommands = []string{"echo", "exit", "type"}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		command, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		command = strings.TrimRight(command, "\n")
		if command == "exit 0" {
			os.Exit(0)
		} else if strings.HasPrefix(command, "echo ") {
			fmt.Println(command[5:])
		} else if strings.HasPrefix(command, "type ") {
			if slices.Contains(validCommands, command[5:]) {
				fmt.Printf("%s is a shell builtin\n", command[5:])
			} else {
				fmt.Printf("%s: not found\n", command[5:])
			}
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
