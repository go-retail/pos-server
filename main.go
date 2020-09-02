package main

import (
	"log"
	"net/http"

	"github.com/go-retail/common-utils/pkg/configutils"
	"github.com/go-retail/common-utils/pkg/rabbit"
	routes "github.com/go-retail/pos-server/pkg/routes"
)

func main() {
	configutils.GetConfig()
	rabbit.InitRMQ()
	defer rabbit.Rmq.Connection.Close()
	defer rabbit.Rmq.Channel.Close()

	router := routes.NewRouter()
	log.Println("Starting Service on port :8080 ")
	log.Fatal(http.ListenAndServe(":8080", router))

}
