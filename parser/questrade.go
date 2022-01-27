package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type QuestradeParser struct{}

func (qparser QuestradeParser) ParseString(csvString string) (map[string]int, map[string][]Transaction, error) {
	rows := strings.Split(csvString, "\n")

	stocksList := map[string]struct{}{}
	transactionsByDate := make(map[string][]Transaction)
	uniqueSymbols := make(map[string]int)

	// first row should be the columns
	for _, row := range rows[1:] {
		if row == "" {
			continue
		}

		cols := strings.Split(row, ",")

		transactionDate := cols[0]
		action := cols[2]
		symbol := cols[3]
		quantity, err := strconv.ParseFloat(cols[5], 64)
		if err != nil {
			continue
		}
		price, err := strconv.ParseFloat(cols[6], 64)
		if err != nil {
			continue
		}
		gross, err := strconv.ParseFloat(cols[7], 64)
		if err != nil {
			continue
		}

		dateValue, err := time.Parse("2006-01-02 15:04:05 AM", transactionDate)
		dateStr := dateValue.Format("2006-01-02")

		if err != nil {
			fmt.Println(err)
			continue
		}

		parsedAction, err := ParseAction(action)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if _, ok := uniqueSymbols[symbol]; !ok {
			uniqueSymbols[symbol] = 1
		}

		transaction := Transaction{
			Symbol:          symbol,
			Action:          parsedAction,
			Gross:           gross,
			Quantity:        quantity,
			Price:           price,
			TransactionDate: dateValue,
		}

		if value, found := transactionsByDate[dateStr]; found {
			transactionsByDate[dateStr] = append(value, transaction)
		} else {
			transactionsByDate[dateStr] = []Transaction{transaction}
		}

		// add stock to unique stock list
		if _, ok := stocksList[symbol]; !ok {
			stocksList[symbol] = struct{}{}
		}
	}

	return uniqueSymbols, transactionsByDate, nil
}
