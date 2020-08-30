package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type Transaction struct {
	CustFirstName string    `json:"custFirstName"`
	CustLastName  string    `json:"custLastName"`
	Total         float64   `json:"total"`
	TxnDate       time.Time `json:"txnDate"`
}

func main() {
	getConfig()
	startMux()
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
	w.Write([]byte("Hello from Snippetbox"))
}

func startMux() {

	http.HandleFunc("/", home)
	// http.HandleFunc("/snippet", showSnippet)
	http.HandleFunc("/txns/create", createTxn)
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", nil)
	log.Fatal(err)

}

func createTxn(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405)
		return
	}

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
	conn, channel, queue := getQueue()

	defer conn.Close()
	defer channel.Close()

	msgString, _ := json.Marshal(*txn)

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        msgString,
	}

	channel.Publish("", queue.Name, false, false, msg)

}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	//TODO externalize RabbitMQ Port to config
	urlString := fmt.Sprintf("amqp://%s:%s@%s:5672", viper.GetString("username"), viper.GetString("password"), viper.GetString("host"))

	conn, err := amqp.Dial(urlString)
	failOnError(err, "Failed to Connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to Open The Channel")

	//TODO Deliver to Exchange not a Queue
	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	failOnError(err, "Unable to Declare a Queue")

	return conn, ch, &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", err, msg)
		panic(fmt.Sprintf("%s: %s", err, msg))
	}
}

func getConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("conf")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			failOnError(err, "Error reading Config file")
		}
	}
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())
}
