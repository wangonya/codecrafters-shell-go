package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var builtins = []string{"echo", "exit", "type"}

// commandExistsInPath checks if the command exists in any of the locations in $PATH
//
// If a match is found, the path location is returned, otherwise err
func commandExistsInPath(command string) (string, error) {
	for _, path := range strings.Split(os.Getenv("PATH"), ":") {
		if _, err := os.Stat(fmt.Sprintf("%s/%s", path, command)); err != nil {
			continue
		}
		return path, nil
	}
	return "", fmt.Errorf("%s: not found", command)
}

func runCmd(command string) (string, error) {
	split := strings.Split(command, " ")

	_, err := commandExistsInPath(split[0])
	if err != nil {
		return "", err
	}

	out, err := exec.Command(split[0], split[1:]...).Output()
	return string(out), err
}

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
			if slices.Contains(builtins, command[5:]) {
				fmt.Println(command[5:], "is a shell builtin")
				continue
			}
			path, err := commandExistsInPath(command[5:])
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			} else {
				fmt.Printf("%s is %s\n", command[5:], fmt.Sprintf("%s/%s", path, command[5:]))
			}
		} else {
			out, err := runCmd(command)
			if err != nil {
				fmt.Printf("%s: command not found\n", command)
			} else {
				fmt.Print(out)
			}
		}
	}
}
