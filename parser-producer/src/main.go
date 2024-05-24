package main

import (
	"log"
	"os"

	model "github.com/ajandera/sp_model"
	"github.com/ajandera/sp_model/rdbsClientInfo"
	"github.com/bitly/go-simplejson"
	_ "github.com/streadway/amqp"
	"github.com/wagslane/go-rabbitmq"
)

var stores []rdbsClientInfo.Stores

func produce(m model.Repository, publisher *rabbitmq.Publisher) {
	log.Println("Start parser producer")
	// get all accounts
	accounts := m.GetAccounts()
	for _, account := range accounts {
		stores = m.GetStoresByAccount(account.Id.String())
		for _, store := range stores {
			log.Println("Open data parser for store id: " + store.Id.String())
			message := simplejson.New()
			message.Set("storeId", store.Id)
			res, _ := message.MarshalJSON()
			err := publisher.Publish(
				res,
				[]string{"parser"},
				rabbitmq.WithPublishOptionsContentType("application/json"),
			)
			if err != nil {
				log.Printf(err.Error())
				return
			}
		}
	}
}

func main() {
	clientsDataDsn := "host=" + os.Getenv("CLIENTS_DATA_HOST") + " user=" + os.Getenv("CLIENTS_DATA_USER") + " password=" + os.Getenv("CLIENTS_DATA_PASSWORD") + " dbname=" + os.Getenv("CLIENTS_DATA_DATABASE") + " port=" + os.Getenv("CLIENTS_DATA_PORT") + " sslmode=disable"
	clientsInformationDataDsn := "host=" + os.Getenv("CLIENTS_INFORMATION_HOST") + " user=" + os.Getenv("CLIENTS_INFORMATION_USER") + " password=" + os.Getenv("CLIENTS_INFORMATION_PASSWORD") + " dbname=" + os.Getenv("CLIENTS_INFORMATION_DATABASE") + " port=" + os.Getenv("CLIENTS_INFORMATION_PORT") + " sslmode=disable"
	repository := model.ClientsInit(
		clientsDataDsn,
		clientsInformationDataDsn)

	log.Println("Connect to databases")
	log.Println("Start publishing")
	conn, err := rabbitmq.NewConn(
		"amqp://"+os.Getenv("RABBIT_USER")+":"+os.Getenv("RABBIT_PASS")+"@"+os.Getenv("RABBIT_HOST"),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	publisher, err := rabbitmq.NewPublisher(conn)
	if err != nil {
		log.Println(err)
	}
	defer publisher.Close()

	produce(repository, publisher)
}
