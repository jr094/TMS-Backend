package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type JsonTransaction struct {
	Symbol          string  `json:"Symbol"`
	Action          string  `json:"Action"`
	Gross           float64 `json:"Gross"`
	Quantity        float64 `json:"Quantity"`
	Price           float64 `json:"Price"`
	TransactionDate string  `json:"TransactionDate"`
}

type LocalDataParser struct{}

func (localparser LocalDataParser) ParseBrokerData() (map[string][]Transaction, error) {
	jsonFile, err := os.Open("mydata.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}

	var transactions []JsonTransaction
	json.Unmarshal(data, &transactions)

	transactionsByDate := make(map[string][]Transaction)
	for _, t := range transactions {
		dateStr := t.TransactionDate
		timeParsedDate, err := time.Parse("2006-01-02", dateStr)

		parsedAction, err := ParseAction(t.Action)
		if err != nil {
			fmt.Println(err)
		}

		transaction := Transaction{
			Symbol:          t.Symbol,
			Action:          parsedAction,
			Gross:           t.Gross,
			Quantity:        t.Quantity,
			Price:           t.Price,
			TransactionDate: timeParsedDate,
		}

		if value, found := transactionsByDate[dateStr]; found {
			transactionsByDate[dateStr] = append(value, transaction)
		} else {
			transactionsByDate[dateStr] = []Transaction{transaction}
		}
	}

	return transactionsByDate, nil
}
