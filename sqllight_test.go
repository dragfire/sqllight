package main

import (
    "log"
    "strings"
    "testing"
    "os/exec"
)

func fatalErrCheck(err error) {
    if err != nil {
	log.Fatal(err)
    }
}

func equal(got, want []string) bool {
    gotString := strings.Join(got, "")
    wantString := strings.Join(got, "")

    //fmt.Printf("%s\n%s\n", gotString, wantString)
    return strings.Compare(gotString, wantString) == 0
}

func runDB(cmds []string) []string  {
    cmd := exec.Command("./sqllight")

    stdin, err := cmd.StdinPipe()
    fatalErrCheck(err)

    go func() {
	defer stdin.Close()
	stdin.Write([]byte(strings.Join(cmds[:],"\n")))
    }()

    out, err := cmd.Output()
    fatalErrCheck(err)

    //fmt.Printf("%s\n\n", string(out))
    return strings.Split(string(out), "\n")
}

func TestSqllight(t *testing.T) {
    t.Run("insert and retrieves a row", func(t *testing.T) {
	got := runDB([]string{"insert 1 a b", "select", ".exit"})
	want := []string{"sqllight > Executed", "sqllight > User: 1 a b", "Executed", "sqllight > "}

	if !equal(got, want) {
	    t.Errorf("got: %v, want: %v", got, want)
	}
    })
}
