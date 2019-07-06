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
    PrepareSyntaxError
    PrepareUnrecognizedStatement
)

type StatementType int

const (
    StatementTypeInsert StatementType = iota
    StatementTypeSelect
)

type ExecuteResult int

const (
    ExecuteSuccess ExecuteResult = iota
    ExecuteTableFull
)

const (
    ColumnUsernameSize int = 32
    ColumnEmailSize    int = 255
)

type Row struct {
    id       uint32
    username [ColumnUsernameSize]rune
    email    [ColumnEmailSize]rune
}

type Statement struct {
    statementType StatementType
    rowToInsert   *Row
}

const (
    TableMaxPages = 100
    RowsPerPage = 200
)

type Page struct {
    rows [RowsPerPage]*Row
}

type Table struct {
    numRows uint32
    pages   [TableMaxPages]*Page
}

func newTable() *Table {
    table := &Table{numRows: 0}
    for i := 0; i < TableMaxPages; i++ {
	table.pages[i] = nil
    }
    return table
}

func rowSlot(table *Table, rowNum uint32) {
    pageNum := rowNum / RowsPerPage
    page := table.pages[pageNum]

    if page == nil {
	page = new(Page)
    }
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
    if strings.LastIndex(cmd, "insert") == 0 {
	statement.statementType = StatementTypeInsert
	rowToInsert := statement.rowToInsert

	var id uint32
	var username, email string

	argsAssigned, err := fmt.Sscanf(cmd, "insert %d %s %s", &id, &username, &email)
	if argsAssigned < 3 || err != nil {
	    fmt.Println(err)
	    return PrepareSyntaxError
	}

	rowToInsert.id = id
	copy(rowToInsert.username[:], []rune(username))
	copy(rowToInsert.email[:], []rune(email))

	return PrepareSuccess
    }

    if strings.LastIndex(cmd, "select") == 0 {
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
    table := newTable()
    fmt.Printf("%q\n", table)
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

	statement := Statement{rowToInsert: new(Row)}
	switch prepareStatement(cmd, &statement) {
	case PrepareSuccess:
	    break
	case PrepareSyntaxError:
	    fmt.Println("Syntax error. Could not parse statement")
	    continue
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
