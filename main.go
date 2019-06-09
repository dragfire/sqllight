package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	cmdExit = ".exit"
	cmdHelp = ".help"
)

func displayLeftDbName() {
	fmt.Print("sqllight > ")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	displayLeftDbName()
	for scanner.Scan() {
		cmd := scanner.Text()
		switch cmd {
		case cmdExit:
			os.Exit(0)
		default:
			fmt.Println("Unknown command: ", cmd)
		}
		displayLeftDbName()
	}

	if scanner.Err() != nil {
		log.Panic(scanner.Err())
	}
}
