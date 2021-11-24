package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	Symbol          string
	Action          string
	Gross           float64
	Quantity        float64
	Price           float64
	TransactionDate string
}

func ParseQuestTrade(csvString string) (map[string][]Transaction, error) {
	rows := strings.Split(csvString, "\n")

	transactionsByDate := make(map[string][]Transaction)

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
		price, err := strconv.ParseFloat(cols[6], 64)
		gross, err := strconv.ParseFloat(cols[7], 64)

		dateValue, err := time.Parse("2006-01-02 15:04:05 AM", transactionDate)
		dateStr := dateValue.Format("2006-01-02")

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		transaction := Transaction{
			Symbol:          symbol,
			Action:          action,
			Gross:           gross,
			Quantity:        quantity,
			Price:           price,
			TransactionDate: dateStr,
		}

		if value, found := transactionsByDate[dateStr]; found {
			transactionsByDate[dateStr] = append(value, transaction)
		} else {
			transactionsByDate[dateStr] = []Transaction{transaction}
		}
	}

	return transactionsByDate, nil
}
