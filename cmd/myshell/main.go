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
		fmt.Printf("%s: command not found\n", strings.TrimRight(response, "\n"))
	}
}
