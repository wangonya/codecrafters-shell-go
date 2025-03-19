package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

var builtins = []string{"echo", "exit", "type", "pwd"}

type command struct {
	executable string
	args       []string
}

// commandExistsInPath checks if the command exists in any of the locations in $PATH
//
// If a match is found, the path location is returned, otherwise err
func commandExistsInPath(command string) (string, error) {
	if strings.HasPrefix(command, "'") {
		command = strings.Trim(command, "'")
	} else if strings.HasPrefix(command, "\"") {
		command = strings.Trim(command, "\"")
	}
	for _, path := range strings.Split(os.Getenv("PATH"), ":") {
		if _, err := os.Stat(fmt.Sprintf("%s/%s", path, command)); err != nil {
			continue
		}
		return path, nil
	}
	return "", fmt.Errorf("%s: not found", command)
}

func runCmd(cmd command) string {
	_, err := commandExistsInPath(cmd.executable)
	if err != nil {
		return err.Error()
	}

	out, err := exec.Command("sh", "-c", strings.Join(append([]string{cmd.executable}, cmd.args...), " ")).CombinedOutput()
	return string(out)
}

// filterEmptyArgs removes any leading and trailing spaces that
// may be inteprated as args
func filterEmptyArgs(args []string) []string {
	x := []string{}
	inQuotes := false
	for _, v := range args {
		if strings.HasPrefix(v, "'") || strings.HasPrefix(v, "\"") {
			inQuotes = true
		} else if inQuotes && (strings.HasSuffix(v, "'") || strings.HasSuffix(v, "\"")) {
			inQuotes = false
		}

		if v != "" || inQuotes {
			x = append(x, v)
		}
	}
	return x
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}

		splitInput := strings.Split(strings.TrimRight(input, "\n"), " ")
		splitInput = filterEmptyArgs(splitInput)
		cmd := command{}
		if strings.HasPrefix(splitInput[0], "'") || strings.HasPrefix(splitInput[0], "\"") {
			prefix := string(splitInput[0][0])
			x := []string{splitInput[0]}
			for _, v := range splitInput[1:] {
				x = append(x, v)
				if strings.HasSuffix(v, prefix) {
					break
				}
			}
			cmd.executable = strings.Join(x, " ")
			cmd.args = splitInput[len(x):]
		} else {
			cmd.executable = splitInput[0]
			cmd.args = splitInput[1:]
		}

		switch cmd.executable {
		case "echo":
			out := runCmd(cmd)
			fmt.Print(out)
		case "exit":
			exitCode, err := strconv.Atoi(cmd.args[0])
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(exitCode)
		case "type":
			if slices.Contains(builtins, cmd.args[0]) {
				fmt.Println(cmd.args[0], "is a shell builtin")
				continue
			}
			path, err := commandExistsInPath(cmd.args[0])
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			} else {
				fmt.Printf("%s is %s\n", cmd.args[0], fmt.Sprintf("%s/%s", path, cmd.args[0]))
			}
		case "pwd":
			pwd, _ := os.Getwd()
			fmt.Println(pwd)
		case "cd":
			path := cmd.args[0]
			if path == "~" {
				path, _ = os.UserHomeDir()
			}
			err := os.Chdir(path)
			if err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", path)
			}
		case "cat":
			out := runCmd(cmd)
			fmt.Print(out)
		default:
			out := runCmd(cmd)
			fmt.Print(out)
		}
	}
}
