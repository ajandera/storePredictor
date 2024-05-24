package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gocolly/colly/v2"
	_ "github.com/streadway/amqp"
	"github.com/wagslane/go-rabbitmq"
	"gitlab.eaineu.com/storepredictor/model"
	"gitlab.eaineu.com/storepredictor/model/rdbsClientInfo"
)

const consumerName = "parser"

type CrawledData struct {
	Index      int
	ParsedData float64
}

type Message struct {
	StoreId string
}

func parser(store rdbsClientInfo.Stores, repository model.Repository) {
	webalizeUrl := strings.Replace(store.Url, ".", "-", -1)
	baseUrl := "https://obchody.heureka.cz/" + webalizeUrl + "/recenze/overene"
	c := colly.NewCollector()
	data := []CrawledData{}

	c.OnHTML(".c-shop-detail-recommendation__percentage,c-stats-table__badge", func(e *colly.HTMLElement) {
		item := CrawledData{}
		val, err := strconv.ParseFloat(e.Text, 64)
		if err != nil {
			log.Println(err)
		}
		item.Index = e.Index
		item.ParsedData = val
		data = append(data, item)
	})

	c.OnScraped(func(r *colly.Response) {
		var customerSatisfaction float64
		var storePower float64
		for _, cd := range data {
			if cd.Index == 0 {
				customerSatisfaction = cd.ParsedData
			} else if cd.Index == 1 {
				storePower = cd.ParsedData
			}
		}
		if storePower <= 0 {
			storePower = customerSatisfaction * store.PerceivedValue * 1.2
		}
		repository.CreateOpenData(storePower, customerSatisfaction, store.MinimalProductPrice,
			store.MaximalProductPrice, store.PerceivedValue, store.Id.String())
	})

	c.Visit(baseUrl)
	update(repository)
}

func update(m model.Repository) {
	log.Println("Start parameter updater")

	// get all accounts
	accounts := m.GetAccounts()
	for _, account := range accounts {
		// get store for account
		stores := m.GetStoresByAccount(account.Id.String())
		for _, store := range stores {
			// get opendata for store
			log.Println(store)
			// save last open data information

		}
	}
}

func main() {
	clientsDataDsn := "host=" + os.Getenv("CLIENTS_DATA_HOST") + " user=" + os.Getenv("CLIENTS_DATA_USER") + " password=" + os.Getenv("CLIENTS_DATA_PASSWORD") + " dbname=" + os.Getenv("CLIENTS_DATA_DATABASE") + " port=" + os.Getenv("CLIENTS_DATA_PORT") + " sslmode=disable"
	clientsInformationDataDsn := "host=" + os.Getenv("CLIENTS_INFORMATION_HOST") + " user=" + os.Getenv("CLIENTS_INFORMATION_USER") + " password=" + os.Getenv("CLIENTS_INFORMATION_PASSWORD") + " dbname=" + os.Getenv("CLIENTS_INFORMATION_DATABASE") + " port=" + os.Getenv("CLIENTS_INFORMATION_PORT") + " sslmode=disable"
	repository := model.ClientsInit(
		clientsDataDsn,
		clientsInformationDataDsn)

	conn, err := rabbitmq.NewConn(
		"amqp://"+os.Getenv("RABBIT_USER")+":"+os.Getenv("RABBIT_PASS")+"@"+os.Getenv("RABBIT_HOST"),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	// wait for server to acknowledge the cancel
	const noWait = false
	var message Message
	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("consumed: %v", string(d.Body))
			errMarsh := json.Unmarshal(d.Body, &message)

			if errMarsh != nil {
				fmt.Println(errMarsh)
			}
			parser(repository.GetStoreById(message.StoreId), repository)
			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.Ack
		},
		"parser",
		rabbitmq.WithConsumerOptionsConcurrency(1),
		rabbitmq.WithConsumerOptionsConsumerName(consumerName),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	// block main thread - wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("stopping consumer")
}
