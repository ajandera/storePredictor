package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/models/cbhmm"
	"main/models/lp"
	"main/utils"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gonum/matrix/mat64"
	"github.com/montanaflynn/stats"
	_ "github.com/streadway/amqp"
	"github.com/wagslane/go-rabbitmq"
	model "gitlab.eaineu.com/storepredictor/model"
	"gitlab.eaineu.com/storepredictor/model/rdbsClientData"
	"gitlab.eaineu.com/storepredictor/model/rdbsClientInfo"
)

const cbhmmLevel = 3
const predictionConstant = 3
const consumerName = "predictor"
const All = 0
const Cbhmm = 1
const Products = 2

var states = []string{"order", "order not completed", "no order"}

// Message struct for RabbitMQ
type Message struct {
	StoreId string
	D0      string
	Product string
	Type    int
	Plan    bool
}

// calculate function to predict data
func calculate(m model.Repository, i model.Influx, storeId string, dayString string, typeOfPrediction int, productDirect string, plan bool) {
	log.Println("Start prediction")
	// define variables
	var predictedOrders []float64
	var visitors stats.Float64Data
	var actualVisitors []rdbsClientData.VisitorsByDay
	var preWeights []float64
	var pctWages []float64
	var preWeightsOrders []float64
	var pctWagesOrder []float64
	var visit int = 0

	// parse time from rabbit
	t, dayTimeErr := time.Parse(time.RFC3339, dayString)

	if dayTimeErr != nil {
		log.Println(dayTimeErr.Error())
	}

	// get store object
	store := m.GetStoreById(storeId)
	predictedPeriod := store.Window

	// get store weights
	storeWeights := m.GetStoreWeights(store.Id.String())
	// FIx for speed and optimization, get only last month
	// get first day of tracking
	//tFrom := m.GetFirstRecord(map[string]interface{}{"store_id": store.Id})
	later := t.Add(-24 * 31 * time.Hour)
	tFrom := time.Date(later.Year(), later.Month(), later.Day(), 0, 0, 0, 0, later.Location()).Format(time.RFC3339)
	// get actual visitors for set period
	actualVisitors = m.GetVisitorsForPredictionView(tFrom, dayString, store.Id.String())
	// get all orders from first track
	previousOrders := int(m.GetSumOrdersForPrediction(
		map[string]interface{}{"store_id": store.Id, "created": dayString}))

	// check for new stores
	if previousOrders == 0 {
		previousOrders = 5000
	}

	// predict when visitors exist
	if len(actualVisitors) >= 7 {
		log.Println("Predict store id: " + store.Id.String())
		visitors = make(stats.Float64Data, 0)
		for _, v := range actualVisitors {
			visitors = append(visitors, float64(v.Visitors))
			visit += v.Visitors
		}
		if visit > 0 {
			// calculate store margin
			margin := float64((previousOrders / visit) * 100)

			// calculate statistic variables for dataset to get year data
			median, _ := stats.Median(visitors)
			std, _ := stats.StandardDeviation(visitors)

			// diff each day visitors from statistic variables
			for _, visitor := range visitors {
				preWeight := visitor - (median - (std / float64(predictionConstant)))
				preWeights = append(preWeights, preWeight)
				// percentage diff
				pctWages = append(pctWages, 1-((visitor-preWeight)/visitor))
			}

			pctWage, _ := stats.Median(pctWages)
			predictVisitors := lp.Elp(predictedPeriod, predictionConstant, visitors, pctWage)
			log.Println("Start turn over prediction")
			if plan == false {
				log.Println("Predict store id: " + store.Id.String() + " with plan false")
				// coefficient data normally get google analytics
				// normal store from [0.75, 0.95]
				perceived := store.PerceivedValue

				// Number of store products
				uniqueProductsSell := store.ProductSell

				// Get from open data, usually from [0,85, 0,99]
				customerSatisfaction := store.ActualCustomerSatisfaction

				// number of product order each day
				// this is initial value and it's dynamically calculate form predicted values
				// in the model prediction loop
				averageQ := float64(previousOrders / uniqueProductsSell)
				averageP := margin

				// Quality index for vendor model
				// default from open data, usually from [0.8, 0.95]
				// 1 for market domination
				Q1 := store.ActualStorePower

				// coefficients for vendor probability
				beta := storeWeights.Beta
				gamma := storeWeights.Gama
				delta := storeWeights.Delta

				// loyalty model coefficients
				// industry and country dependent
				// can be set from eg. google analytics data analysis
				a := storeWeights.A
				b := storeWeights.B
				c := storeWeights.C
				d := storeWeights.D
				e := storeWeights.E

				// lower product price
				minimalProductSell := store.MinimalProductPrice

				// higher product price
				maximalProductSell := store.MaximalProductPrice

				experimentationCycles := 2

				// Probability matrix is used for wages to our sub-model
				emissionMatrix := utils.Matrix(storeWeights.ProbabilityWeights)

				// prediction cycle
				if typeOfPrediction == All || typeOfPrediction == Cbhmm {
					for _, dayVisitors := range predictVisitors {
						orders := 0
						if dayVisitors > 0 {
							for i := 1; i < experimentationCycles; i++ {
								trust := (1 / previousOrders) * previousOrders / int(dayVisitors)
								for pu := 1; pu < int(dayVisitors); pu++ {
									Pj := utils.GeneratePriceOfProduct(minimalProductSell, maximalProductSell)
									// product price minus product retail price
									Pi := Pj * margin
									// actual product order value for each cycle
									Qj := 1

									randValue := rand.Float64()
									var Q2 float64
									var P2 float64
									var satisfaction float64

									// Quality index for vendor 2
									min := 70
									max := 100
									Q2 = float64((rand.Intn(max-min) + min) / 100)

									// vendor 2 product price
									min2 := Pj - (Pj / 20)
									max2 := Pj + (Pj / 20)
									P2 = float64(rand.Intn(int(max2)-int(min2)) + int(min2))
									satisfaction = customerSatisfaction - randValue/10

									// get value of vendor sub-model
									vendorProbability := cbhmm.Vendor(beta, gamma, delta, Pj, P2, Q1, Q2, store.Id.String())

									// get value of psychology sub-model
									psychologyProbability := cbhmm.Psychology(averageQ, averageP, Pj, Pi, float64(Qj), float64(orders))

									// get value of loyalty sub-model
									loyaltyProbability := cbhmm.Loyalty(satisfaction, float64(trust), perceived, a, b, c, d, e, store.Id.String())
									if vendorProbability == 0 {
										continue
									}

									var transitionMatrix mat64.Matrix
									transitionMatrixString := "[" + fmt.Sprintf("%f", vendorProbability) + " " + fmt.Sprintf("%f", psychologyProbability) + " " + fmt.Sprintf("%f", loyaltyProbability) + ";" + fmt.Sprintf("%f", vendorProbability) + " " + fmt.Sprintf("%f", psychologyProbability) + " " + fmt.Sprintf("%f", loyaltyProbability) + ";" + fmt.Sprintf("%f", vendorProbability) + " " + fmt.Sprintf("%f", psychologyProbability) + " " + fmt.Sprintf("%f", loyaltyProbability) + "]"
									transitionMatrix = utils.Matrix(transitionMatrixString)
									path, _, last := cbhmm.Cbhmm(cbhmmLevel, transitionMatrix, emissionMatrix, states)
									S := make([]int, 3)
									for ii := 0; ii < len(path)-1; ii++ {
										S[ii] = utils.IndexOf(path[ii], states)
									}
									read := cbhmm.Viterbi(S, S, transitionMatrix, emissionMatrix)

									if last == states[0] && utils.Sum(read[1:3]) < 2 {
										orders++
										// increase psychology effect of success order in store
										Qj++
									}
								}
							}
						}
						predictedOrders = append(predictedOrders, math.Ceil(float64(orders)/30))
					}
				}
			} else {
				log.Println("Predict store id: " + store.Id.String() + " with plan true")
				var orders []float64
				previousOrders := m.GetOrdersForPrediction(tFrom, dayString, store.Id.String())
				for _, o := range previousOrders {
					orders = append(visitors, float64(o.Orders))
				}
				// calculate statistic variables for dataset to get year data
				medianOrder, _ := stats.Median(visitors)
				stdOrder, _ := stats.StandardDeviation(visitors)

				// diff each day visitors from statistic variables
				for _, order := range orders {
					preWeightOrder := order - (medianOrder - (stdOrder / float64(predictionConstant)))
					preWeightsOrders = append(preWeightsOrders, preWeightOrder)
					// percentage diff
					pctWagesOrder = append(pctWagesOrder, 1-((order-preWeightOrder)/order))
				}

				pctWageOrder, _ := stats.Median(pctWagesOrder)
				predictedOrders = lp.Elp(predictedPeriod, predictionConstant, visitors, pctWageOrder)
			}

			avg := m.GetAvgAmountForPrediction(map[string]interface{}{"store_id": store.Id})

			// save order d0
			td0o := t
			day := time.Date(td0o.Year(), td0o.Month(), td0o.Day(), 0, 0, 0, 0, td0o.Location())
			from := time.Date(td0o.Year(), td0o.Month(), td0o.Day(), 0, 0, 0, 0, td0o.Location()).Format(time.RFC3339)
			to := time.Date(td0o.Year(), td0o.Month(), td0o.Day(), 23, 59, 59, 0, td0o.Location()).Format(time.RFC3339)
			_, err := i.StoreData("order", "d0", int(m.GetOrdersCountByDate(map[string]interface{}{"from": from, "to": to, "store_id": store.Id})), avg, day, store.Id.String(), os.Getenv("INFLUX_ORGANIZATION"))

			if err != nil {
				log.Println(err.Error())

			}

			// save orders
			for or := 1; or <= len(predictedOrders); or++ {
				t := t.AddDate(0, 0, or)
				day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
				_, err := i.StoreData("order", fmt.Sprintf("%s%02d", "d", or), int(predictedOrders[or-1]), avg, day, store.Id.String(), os.Getenv("INFLUX_ORGANIZATION"))
				if err != nil {
					log.Println(err.Error())

				}
			}
			// save visitors
			// add d0 visitors
			_, errV := i.StoreData("visitors", "d0", int(m.GetVisitorsCountByDate(map[string]interface{}{"from": from, "to": to, "store_id": store.Id})), 0, day, store.Id.String(), os.Getenv("INFLUX_ORGANIZATION"))
			if errV != nil {
				log.Println(errV.Error())

			}
			// add prediction
			for vr := 1; vr <= len(predictVisitors); vr++ {
				t := t.AddDate(0, 0, vr)
				day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
				_, err := i.StoreData("visitors", fmt.Sprintf("%s%02d", "d", vr), int(predictVisitors[vr-1]), 0, day, store.Id.String(), os.Getenv("INFLUX_ORGANIZATION"))
				if err != nil {
					log.Println(err.Error())

				}
			}
			log.Println("Stop global prediction")
			if (typeOfPrediction == All || typeOfPrediction == Products) && plan == false {
				log.Println("Start product prediction")
				products := m.GetProducts(map[string]interface{}{"store_id": store.Id, "limit": 10000, "offset": 0})
				for _, p := range products {
					if productDirect != "" {
						if productDirect == p.ProductCode {
							productPrediction(m, i, p.ProductCode, tFrom, t.Format(time.RFC3339), store)
						}
					} else {
						productPrediction(m, i, p.ProductCode, tFrom, t.Format(time.RFC3339), store)
					}
				}
				log.Println("Stop product prediction")
			}
		}
		_, err := i.Flush(os.Getenv("INFLUX_ORGANIZATION"), store.Id.String())
		if err != nil {

			return
		}
	} else {
		log.Println("Skip store id: " + store.Id.String() + " not enough data")
	}
	log.Println("End prediction")
}

// productPrediction function to predict per product
func productPrediction(m model.Repository, i model.Influx, productCode string, from string, to string, store rdbsClientInfo.Stores) {
	var visitorsVector []float64
	var ordersVector []float64
	var pctWagesP []float64
	var d0 int
	visitValue := 0
	orderValue := 0

	// get visitors vector
	visitorsVectorDb := m.GetVisitorsForPredictionPerProductView(from, to, store.Id.String(), productCode)
	// get orders vector
	ordersVectorDb := m.GetOrdersForPredictionPerProductView(from, to, store.Id.String(), productCode)

	// prediction visitors for product
	for _, v := range visitorsVectorDb {
		visitorsVector = append(visitorsVector, float64(v.Visitors))
		visitValue += v.Visitors
	}

	if len(visitorsVector) > 10 && visitValue > 0 {

		// prediction orders for product
		for _, o := range ordersVectorDb {
			ordersVector = append(ordersVector, float64(o.Orders))
			orderValue += o.Orders
		}

		//marginPerProduct := float64(orderValue) / float64(visitValue)

		// calculate statistic variables for dataset to get year data
		medianP, _ := stats.Median(visitorsVector)
		stdP, _ := stats.StandardDeviation(visitorsVector)

		// diff each day visitors from statistic variables
		for _, visitorP := range visitorsVector {
			if visitorP > 0 {
				preWeightP := medianP - stdP
				// percentage diff
				pctWagesP = append(pctWagesP, 1-((visitorP-preWeightP)/visitorP))
			} else {
				pctWagesP = append(pctWagesP, 0)
			}
		}
		// get mean value
		pctWageP, _ := stats.Mean(pctWagesP)
		orderPredictionP := lp.Elp(store.Window, predictionConstant, visitorsVector, pctWageP)

		// save to influx
		// save order d0
		tpD0o, _ := time.Parse(time.RFC3339, to)
		day := time.Date(tpD0o.Year(), tpD0o.Month(), tpD0o.Day(), 0, 0, 0, 0, tpD0o.Location())
		fromP := day.Format(time.RFC3339)
		toP := time.Date(tpD0o.Year(), tpD0o.Month(), tpD0o.Day(), 23, 59, 59, 0, tpD0o.Location()).Format(time.RFC3339)
		d0 = int(m.GetOrdersCountByDatePerProduct(map[string]interface{}{"from": fromP, "to": toP, "store_id": store, "product_code": productCode}))

		_, err := i.StoreData(
			productCode,
			"d0",
			d0,
			0,
			day,
			store.Id.String(),
			os.Getenv("INFLUX_ORGANIZATION"))

		if err != nil {
			log.Println(err.Error())

		}

		// save orders for product
		for or := 1; or <= len(orderPredictionP); or++ {
			t := tpD0o.AddDate(0, 0, or)
			day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			_, err := i.StoreData(
				productCode,
				fmt.Sprintf("%s%02d", "d", or),
				int(orderPredictionP[or-1]),
				0,
				day,
				store.Id.String(),
				os.Getenv("INFLUX_ORGANIZATION"))

			if err != nil {
				log.Println(err.Error())

			}
		}
		saveProductsToStore(m, productCode, store.Id.String(), orderPredictionP, tpD0o)
	}
}

// saveProductsToStore function to product to store
func saveProductsToStore(m model.Repository, productCode string, store string, orderPredictionP []float64, tpD0o time.Time) {
	qty := 0
	c := 1
	for or := 1; or <= len(orderPredictionP); or++ {
		t := tpD0o.AddDate(0, 0, or)
		day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		qty += int(orderPredictionP[or-1])
		if c == 5 {
			dayToNeed := time.Date(t.Year(), t.Month(), t.Day()-3, 0, 0, 0, 0, t.Location())
			record := m.GetProductToStore(productCode, store)
			if utils.IsValidUUID(record.Id) {
				m.UpdateProductToStore(productCode, store, int8(qty), day, dayToNeed)
			} else {
				m.CreateProductToStore(productCode, int8(qty), store, day, dayToNeed)
			}
			qty = 0
			c = 1
		}
		c++
	}
}

// main function to run prediction model
func main() {
	clientsDataDsn := "host=" + os.Getenv("CLIENTS_DATA_HOST") + " user=" + os.Getenv("CLIENTS_DATA_USER") + " password=" + os.Getenv("CLIENTS_DATA_PASSWORD") + " dbname=" + os.Getenv("CLIENTS_DATA_DATABASE") + " port=" + os.Getenv("CLIENTS_DATA_PORT") + " sslmode=disable"
	clientsInformationDataDsn := "host=" + os.Getenv("CLIENTS_INFORMATION_HOST") + " user=" + os.Getenv("CLIENTS_INFORMATION_USER") + " password=" + os.Getenv("CLIENTS_INFORMATION_PASSWORD") + " dbname=" + os.Getenv("CLIENTS_INFORMATION_DATABASE") + " port=" + os.Getenv("CLIENTS_INFORMATION_PORT") + " sslmode=disable"
	repository := model.ClientsInit(
		clientsDataDsn,
		clientsInformationDataDsn)

	influxURL := os.Getenv("INFLUX_HOST")
	influx := model.ClientPredictedDataInit(
		influxURL,
		os.Getenv("INFLUX_TOKEN"))
	log.Println("Connect to databases")

	conn, err := rabbitmq.NewConn(
		"amqp://"+os.Getenv("RABBIT_USER")+":"+os.Getenv("RABBIT_PASS")+"@"+os.Getenv("RABBIT_HOST"),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Println(err)

	}
	defer conn.Close()

	// wait for server to acknowledge the cancel
	var message Message
	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Println("consumed: " + string(d.Body))
			errMarsh := json.Unmarshal(d.Body, &message)

			if errMarsh != nil {
				fmt.Println(errMarsh)

			}
			calculate(repository, influx, message.StoreId, message.D0, message.Type, message.Product, message.Plan)

			// rabbitmq.Ack to send info for continue to next message
			return rabbitmq.Ack
		},
		"prediction",
		rabbitmq.WithConsumerOptionsConcurrency(1),
		rabbitmq.WithConsumerOptionsConsumerName(consumerName),
		rabbitmq.WithConsumerOptionsConsumerAutoAck(true),
	)
	if err != nil {
		log.Println(err)

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
