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
	PrepareNegativeId
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

const (
	TableMaxPages = 10
	RowsPerPage   = 100
	TableMaxRows  = TableMaxPages * RowsPerPage
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

type Page struct {
	rows [RowsPerPage]*Row
}

func newPage() *Page {
}

type Table struct {
	numRows uint32
	pager   *FilePager
}

func dbOpen(filename string) *Table {
	pager := pagerOpen(filename)
	table := &Table{numRows: 0, pager: pager}

	return table
}

func rowSlot(table *Table, row *Row) {
	rowNum := table.numRows
	pageNum := rowNum / RowsPerPage
	page := getPage(table.pager, pageNum)

	if page == nil {
		var nRows [RowsPerPage]*Row
		rows := make([]*Row, RowsPerPage)
		copy(nRows[:], rows)
		table.pages[pageNum] = &Page{nRows}
		page = table.pages[pageNum]
	}

	index := rowNum % RowsPerPage
	page.rows[index] = row
}

func executeInsert(statement *Statement, table *Table) ExecuteResult {
	if table.numRows >= TableMaxRows {
		return ExecuteTableFull
	}

	rowToInsert := statement.rowToInsert
	rowSlot(table, rowToInsert)
	table.numRows++
	return ExecuteSuccess
}

func executeSelect(statement *Statement, table *Table) ExecuteResult {
	var i uint32
	numPages := table.numRows / RowsPerPage
	excessRows := int(table.numRows % RowsPerPage)

	if excessRows > 0 {
		numPages++
	}

	for i = 0; i < numPages; i++ {
		n := RowsPerPage
		rows := table.pages[i].rows
		if i == numPages-1 {
			n = excessRows
		}

		for j := 0; j < n; j++ {
			row := rows[j]
			fmt.Printf("User: %d %s %s\n", row.id, string(row.username[:]), string(row.email[:]))
		}
	}
	return ExecuteSuccess
}

func executeStatement(statement *Statement, table *Table) ExecuteResult {
	switch statement.statementType {
	case StatementTypeInsert:
		return executeInsert(statement, table)
	case StatementTypeSelect:
		return executeSelect(statement, table)
	}
	return ExecuteSuccess
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	table := dbOpen("test.db")
	// fmt.Printf("%q\n", table)
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

		switch executeStatement(&statement, table) {
		case ExecuteSuccess:
			fmt.Println("Executed")
			break
		case ExecuteTableFull:
			fmt.Println("Error: Table full")
			break
		}
	}

	if scanner.Err() != nil {
		log.Panic(scanner.Err())
	}
}
