package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		response, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if response == "exit 0\n" {
			os.Exit(0)
		} else if strings.HasPrefix(response, "echo") {
			fmt.Print(response[5:])
		} else {
			fmt.Printf("%s: command not found\n", strings.TrimRight(response, "\n"))
		}
	}
}
