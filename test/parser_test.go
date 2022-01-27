package parser

import (
	"TMS-Backend/parser"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestQuestradeParser(t *testing.T) {
	curDir, _ := os.Getwd()
	testFile := filepath.Join(curDir, "files/questrade.csv")
	content, err := ioutil.ReadFile(testFile)
	if err != nil || content == nil {
		t.Errorf("Could not parse questrade file - ERR: %s", err)
	}

	questradeParser := parser.QuestradeParser{}
	symbols, transactionsByDate, err := questradeParser.ParseString(string(content))
	if err != nil {
		t.Errorf("Could not parse Questrade CSV string %s", err)
	}

	if _, ok := symbols["AAPL"]; !ok || len(symbols) != 1 {
		t.Errorf("AAPL was not in stocklist after parse")
	}

	expectedCount := 1
	if len(transactionsByDate) != 1 {
		t.Errorf("Transaction count was %d but expected %d", len(transactionsByDate), expectedCount)
	}
}
