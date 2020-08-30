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
	// go client()
	// go server()

	// var s string
	// fmt.Scanln(&s)
	startMux()
}

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

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", home)
	// log.Println("Starting server on :4000")
	// err := http.ListenAndServe(":4000", mux)
	// log.Fatal(err)

	http.HandleFunc("/", home)
	// http.HandleFunc("/snippet", showSnippet)
	http.HandleFunc("/txns/create", createTxn)
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", nil)
	log.Fatal(err)

}

func createTxn(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not. Note that // http.MethodPost is a constant equal to the string "POST".
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// Use the http.Error() function to send a 405 status code and "Method Not // Allowed" string as the response body.
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	//Async Log Transaction Records
	txn := generateTransaction()
	recordTransaction(txn)
	w.Header().Set("Content-Type", "application/json")
	jsonString, err := json.Marshal(&txn)
	w.Write(jsonString)
	failOnError(err, "Unable to Marshal Transaction for HTTP Output")
	// w.Write([]byte(`{"name":"Alex"}`))
	// w.Write([]byte("New Transaction Recorded"))
}

func client() {
	conn, channel, queue := getQueue()

	defer conn.Close()
	defer channel.Close()

	msgs, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register the Channel")

	for msg := range msgs {
		log.Printf("Received message: %s", msg.Body)
	}
}

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
	urlString := fmt.Sprintf("amqp://%s:%s@%s:5672", viper.GetString("username"), viper.GetString("password"), viper.GetString("host"))
	// fmt.Printf(urlString)
	conn, err := amqp.Dial(urlString)
	failOnError(err, "Failed to Connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to Open The Channel")

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
