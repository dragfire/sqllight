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

func prepareStatement(cmd string, statement *Statement) PrepareResult {
    if (strings.LastIndex(cmd, "insert") == 0) {
	statement.statementType = StatementTypeInsert
	return PrepareSuccess
    }
    if (strings.LastIndex(cmd, "select") == 0) {
	statement.statementType = StatementTypeSelect
	return PrepareSuccess
    }

    return PrepareUnrecognizedStatement
}

func executeStatement(statement *Statement) {
    switch statement.statementType {
    case StatementTypeInsert:
	fmt.Println("INSERT OP")
	break
    case StatementTypeSelect:
	fmt.Println("SELECT OP")
	break
    }
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    for {
	displayPrompt()
	scanner.Scan()
	cmd := scanner.Text()
	
	// ignore empty buffer
	if cmd == "" {
	    continue
	}

	if strings.LastIndex(cmd, ".") == 0 {
	    switch performMetaCommand(cmd) {
	    case MetaCommandSuccess:
		break
	    case MetaCommandUnrecognizedCommand:
		fmt.Println("Unknown command: ", cmd)
		continue
	    }
	}

	statement := Statement{}
	switch prepareStatement(cmd, &statement) {
	case PrepareSuccess:
	    break
	case PrepareUnrecognizedStatement:
	    fmt.Printf("Unrecognized keyword at start of '%s'\n", cmd)
	    continue
	}

	executeStatement(&statement)
    }

    if scanner.Err() != nil {
	log.Panic(scanner.Err())
    }
}
