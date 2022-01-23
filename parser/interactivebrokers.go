package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type InteractiveBrokersParser struct{}

func (ibParser InteractiveBrokersParser) ParseString(csvString string) (map[string][]Transaction, map[string]struct{}, error) {
	rows := strings.Split(csvString, "\n")

	stocksList := map[string]struct{}{}
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
		timeParsedDate, err := time.Parse("2006-01-02", transactionDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		action := cols[act_i][1 : len(cols[act_i])-1]
		symbol := cols[sym_i][1 : len(cols[sym_i])-1]
		quantity, err := strconv.ParseInt(cols[q_i][1:len(cols[q_i])-1], 10, 64)
		if err != nil {
			fmt.Println(err)
			continue
		}
		price, err := strconv.ParseFloat(cols[p_i][1:len(cols[p_i])-1], 64)
		if err != nil {
			fmt.Println(err)
			continue
		}
		gross, err := strconv.ParseFloat(cols[g_i][1:len(cols[g_i])-1], 64)
		if err != nil {
			fmt.Println(err)
			continue
		}

		parsedAction, err := ParseAction(action)
		if err != nil {
			fmt.Println(err)
		}

		transaction := Transaction{
			Symbol:          symbol,
			Action:          parsedAction,
			Gross:           gross,
			Quantity:        quantity,
			Price:           price,
			TransactionDate: timeParsedDate,
		}

		if value, found := transactionsByDate[transactionDate]; found {
			transactionsByDate[transactionDate] = append(value, transaction)
		} else {
			transactionsByDate[transactionDate] = []Transaction{transaction}
		}

		// add stock to unique stock list
		if _, ok := stocksList[symbol]; !ok {
			stocksList[symbol] = struct{}{}
		}
	}

	return transactionsByDate, stocksList, nil
}
