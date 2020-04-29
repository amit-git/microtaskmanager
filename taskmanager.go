package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		processOption(presentOptions())
	}
}

func presentOptions() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Pick one >")
	fmt.Println("a <txt> : adds a task ")
	fmt.Println("b <task-id> : puts a task in the backlog")
	fmt.Println("w <task-id> : moves a task out of backlog into current list")
	fmt.Println("d <task-id> : marks task DONE")
	fmt.Println("q : Quit")
	fmt.Printf("\n> ")
	cmd, e := reader.ReadString('\n')
	if e != nil {
		fmt.Println("Error in reading command " + e.Error())
		os.Exit(-1)
	}
	return strings.Replace(cmd, "\n", "", -1)
}

func processOption(cmd string) {
	if cmd == "q" {
		fmt.Println("Ok. Bye")
		os.Exit(0)
	}
	fmt.Println("Processing :: " + cmd)
}
