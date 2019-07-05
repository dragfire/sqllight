package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
)

const (
    cmdExit = ".exit"
    cmdHelp = ".help"
)

type MetaCommandResult int
const (
    MetaCommandSuccess MetaCommandResult = iota
    MetaCommandUnrecognizedCommand
)

type PrepareResult int
const (
    PrepareSuccess PrepareResult = iota
    PrepareUnrecognizedStatement
)

type StatementType int
const (
    StatementTypeInsert StatementType = iota
    StatementTypeSelect
)

type Statement struct {
    statementType StatementType
}

func displayPrompt() {
    fmt.Print("sqllight > ")
}

func performMetaCommand(cmd string) MetaCommandResult {
    if strings.Compare(cmd, ".exit") == 0 {
	os.Exit(0)
    }
    return MetaCommandUnrecognizedCommand
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    for {
	displayPrompt()
	scanner.Scan()
	cmd := scanner.Text()
	if strings.LastIndex(cmd, ".") == 0 {
	    switch performMetaCommand(cmd) {
	    case MetaCommandSuccess:
		break
	    case MetaCommandUnrecognizedCommand:
		fmt.Println("Unknown command: ", cmd)
		break
	    }
	}
    }

    if scanner.Err() != nil {
	log.Panic(scanner.Err())
    }
}
