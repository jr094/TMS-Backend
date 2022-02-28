package parser

import (
	"TMS-Backend/parser"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuestradeParser(t *testing.T) {
	curDir, _ := os.Getwd()
	testFile := filepath.Join(curDir, "files/questrade.csv")
	content, err := ioutil.ReadFile(testFile)

	assert.Nilf(t, err, "Could not parse questrade file - ERR: %s", err)
	assert.NotNilf(t, content, "Could not parse questrade file - ERR: %s", err)

	questradeParser := parser.QuestradeParser{}
	symbols, transactionsByDate, err := questradeParser.ParseString(string(content))
	assert.Nilf(t, err, "Could not parse Questrade CSV string %s", err)

	_, ok := symbols["AAPL"]
	assert.Truef(t, ok, "AAPL was not in stocklist after parse")
	assert.NotEqual(t, 0, len(symbols), "Symbols count was 0")
	assert.NotEqual(t, 0, len(transactionsByDate), "Transaction count was 0")
}

func TestInteractiveParser(t *testing.T) {
	curDir, _ := os.Getwd()
	testFile := filepath.Join(curDir, "files/interactive.csv")
	content, err := ioutil.ReadFile(testFile)

	assert.Nilf(t, err, "Could not parse interactive file - ERR: %s", err)
	assert.NotNilf(t, content, "Could not parse interactive file - ERR: %s", err)

	interactiveBrokersParser := parser.InteractiveBrokersParser{}
	uniqueSymbols, transactionsByDate, err := interactiveBrokersParser.ParseString(string(content))
	assert.Nilf(t, err, "Could not parse interactive CSV string %s", err)

	_, ok := uniqueSymbols["AAPL"]
	assert.Truef(t, ok, "AAPL was not in stocklist after parse")
	assert.Equal(t, 1, len(uniqueSymbols))
	assert.Equal(t, 1, len(transactionsByDate), "Transaction count not 1")
}
