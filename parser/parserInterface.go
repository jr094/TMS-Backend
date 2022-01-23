package parser

import (
	"fmt"
	"strings"
	"time"
)

type Action int64

const (
	Buy Action = iota
	Sell
)

type Transaction struct {
	Symbol          string
	Action          Action
	Gross           float64
	Quantity        int64
	Price           float64
	TransactionDate time.Time
}

type BrokerParser interface {
	ParseString(contents string) (map[string][]Transaction, error)
}

type Document struct {
}

func ParseAction(action string) (Action, error) {
	lowercase := strings.ToLower(action)
	if strings.Contains(lowercase, "buy") || strings.Contains(lowercase, "bought") {
		return Buy, nil
	} else if strings.Contains(lowercase, "sell") || strings.Contains(lowercase, "sold") {
		return Sell, nil
	} else {
		err := fmt.Errorf("could not parse string %s into an Action type", action)
		return -1, err
	}
}

func ParseDocument(filePath string, broker string) (Document, error) {
	return Document{}, nil
}
