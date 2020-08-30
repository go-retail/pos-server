package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

//Transaction .. A Struct to hold a Sale Transaction
type Transaction struct {
	CustFirstName string    `json:"custFirstName"`
	CustLastName  string    `json:"custLastName"`
	Total         float64   `json:"total"`
	TxnDate       time.Time `json:"txnDate"`
}

//This is only for Simulation.  Will not be needed when the Transaction arrives from  POS Machine
func generateTransaction() *Transaction {
	txn :=
		Transaction{
			CustFirstName: "Anand",
			CustLastName:  "Rao",
			Total:         22.34,
			TxnDate:       time.Now(),
		}

	return &txn
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from From the Root"))
}

func createTxn(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	params := mux.Vars(r)
	storeID := params["storeID"]
	posID := params["posID"]
	log.Printf("Received Request from Store: %s -   POS ID: %s", storeID, posID)
	//Generate new Transaction
	//TODO will be removed when real transactions arrive
	txn := generateTransaction()

	//Async Log Transaction Records
	recordTransaction(txn)

	jsonString, err := json.Marshal(&txn)

	if err != nil {
		failOnError(err, "Unable to Marshal Transaction for HTTP Output")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)

}

// func client() {
// 	conn, channel, queue := getQueue()

// 	defer conn.Close()
// 	defer channel.Close()

// 	msgs, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
// 	failOnError(err, "Failed to register the Channel")

// 	for msg := range msgs {
// 		log.Printf("Received message: %s", msg.Body)
// 	}
// }

func recordTransaction(txn *Transaction) {
	msgString, _ := json.Marshal(*txn)

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        msgString,
	}

	rmq.Channel.Publish("", rmq.Queue.Name, false, false, msg)

}
