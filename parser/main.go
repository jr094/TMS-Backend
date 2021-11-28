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

func ParseInteractive(csvString string) (map[string][]Transaction, error) {
	rows := strings.Split(csvString, "\n")

	transactionsByDate := make(map[string][]Transaction)

	// parse first row to get correct column indexes
	var td_i, act_i, sym_i, q_i, p_i, g_i, ac_i int
	headerRow := rows[0]
	for i, col := range strings.Split(headerRow, ",") {
		switch col {
		case "\"TradeDate\"":
			td_i = i
		case "\"Buy/Sell\"":
			act_i = i
		case "\"Symbol\"":
			sym_i = i
		case "\"Quantity\"":
			q_i = i
		case "\"TradePrice\"":
			p_i = i
		case "\"CostBasis\"":
			g_i = i
		case "\"AssetClass\"":
			ac_i = i
		}
	}

	for _, row := range rows[1:] {
		if row == "" {
			continue
		}

		cols := strings.Split(row, ",")
		// go adds extra quotes for some reason
		if cols[ac_i][1:len(cols[ac_i])-1] != "STK" {
			continue
		}

		transactionDate := cols[td_i][1 : len(cols[td_i])-1]
		action := cols[act_i][1 : len(cols[act_i])-1]
		symbol := cols[sym_i][1 : len(cols[sym_i])-1]
		quantity, err := strconv.ParseFloat(cols[q_i][1:len(cols[q_i])-1], 64)
		price, err := strconv.ParseFloat(cols[p_i][1:len(cols[p_i])-1], 64)
		gross, err := strconv.ParseFloat(cols[g_i][1:len(cols[g_i])-1], 64)

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
			TransactionDate: transactionDate,
		}

		if value, found := transactionsByDate[transactionDate]; found {
			transactionsByDate[transactionDate] = append(value, transaction)
		} else {
			transactionsByDate[transactionDate] = []Transaction{transaction}
		}
	}

	return transactionsByDate, nil
}
