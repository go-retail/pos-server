package model

import "time"

//Transaction .. A Struct to hold a Sale Transaction
type Transaction struct {
	CustFirstName string    `json:"custFirstName"`
	CustLastName  string    `json:"custLastName"`
	Total         float64   `json:"total"`
	TxnDate       time.Time `json:"txnDate"`
}
