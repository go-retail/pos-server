package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-retail/pos-server/pkg/rabbit"
	"github.com/go-retail/pos-server/pkg/utils"
	model "github.com/go-retail/retail-model/pkg/model"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

//This is only for Simulation.  Will not be needed when the Transaction arrives from  POS Machine
func generateTransaction() *model.Transaction {
	txn :=
		model.Transaction{
			CustFirstName: "Anand",
			CustLastName:  "Rao",
			Total:         22.34,
			TxnDate:       time.Now(),
		}

	return &txn
}

//Home ..
func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from From the Root"))
}

//CreateTxn ..
func CreateTxn(w http.ResponseWriter, r *http.Request) {

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
		utils.FailOnError(err, "Unable to Marshal Transaction for HTTP Output")
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

func recordTransaction(txn *model.Transaction) {
	msgString, _ := json.Marshal(*txn)

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        msgString,
	}

	rabbit.Rmq.Channel.Publish("", rabbit.Rmq.Queue.Name, false, false, msg)

}
