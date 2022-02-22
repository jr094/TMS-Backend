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

type Broker int64

const (
	UndefinedBroker Broker = iota // Default enum value
	TDAmeritrade
	Questrade
	InteraciveBrokers
)

type Transaction struct {
	Symbol          string
	Action          Action
	Gross           float64
	Quantity        int64
	Price           float64
	TransactionDate time.Time
}

type ParsedDocument struct {
	Symbols      map[string]int
	Transactions map[string][]Transaction
}
type BrokerParser interface {
	ParseString(contents string) (map[string]int, map[string][]Transaction, error)
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

func BrokerParserMap(broker Broker) BrokerParser {
	switch broker {
	case TDAmeritrade:
		return TDAmeritradeParser{}
	case InteraciveBrokers:
		return InteractiveBrokersParser{}
	case Questrade:
		return QuestradeParser{}
	default:
		return nil
	}
}

func ParseDocument(contents string, broker Broker) (ParsedDocument, error) {
	brokerParser := BrokerParserMap(broker)
	if brokerParser == nil {
		err := fmt.Errorf("broker type was undefined")
		return ParsedDocument{}, err
	}
	symbols, transactions, err := brokerParser.ParseString(contents)
	return ParsedDocument{symbols, transactions}, err
}
