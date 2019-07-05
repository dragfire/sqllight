package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "unsafe"
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
    ColumnEmailSize int = 255
)

type Row struct {
    id uint32
    username [ColumnUsernameSize]rune
    email [ColumnEmailSize]rune
}
var row Row

const (
    SizeId uintptr = unsafe.Sizeof(row.id)
    SizeUsername uintptr = unsafe.Sizeof(row.username)
    SizeEmail uintptr = unsafe.Sizeof(row.email)
    SizeRow uintptr = SizeId + SizeUsername + SizeEmail
    IdOffset uintptr = 0
    UsernameOffset uintptr = IdOffset + SizeUsername
    EmailOffset uintptr = SizeUsername + SizeEmail
)


type Statement struct {
    statementType StatementType
    rowToInsert *Row
}

const (
    SizePage uintptr = 4096
    TableMaxPages uintptr = 100
    RowsPerPage = SizePage / SizeRow
    TableMaxRows = RowsPerPage * TableMaxPages
)

type Table struct {
    numRows uint32
    pages [TableMaxPages]interface{}
}

func newTable() *Table {
    table := &Table{numRows: 0}
    var i uintptr = 0
    for ;i<TableMaxPages; i++ {
	table.pages[i] = nil
    }
    return table
}

//func rowSlot(table *Table, uintptr rowNum) uintptr {
//    pageNum := rowNum / RowsPerPage
//    page := table.pages[pageNum]
//
//    if page == nil {
//	page = 
//    }
//}

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
    table := newTable()
    // fmt.Printf("%v\n", table)
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
