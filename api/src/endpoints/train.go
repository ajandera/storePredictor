package endpoints

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	bot "main/bot"
	"main/bot/adapters/logic"
	"main/bot/adapters/storage"
	structs "main/structs/bot"
	utils "main/utils"
	"net/http"
	"os"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

type FeedUrl struct {
	Feed string
}

// POST none
// @tags Train
// @Summary Endpoint to handle training bot
// @Description Endpoint to handle training bot
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Res
// @Router /train/{corpora}  [post]
func ProceedTrain(w http.ResponseWriter, r *http.Request) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	response := simplejson.New()

	// Declare a new TrainData struct.
	var trainData structs.TrainData

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	errBody := json.NewDecoder(r.Body).Decode(&trainData)
	if errBody != nil {
		http.Error(w, errBody.Error(), http.StatusBadRequest)
		log.Printf(errBody.Error())
		return
	}

	vars := mux.Vars(r)
	corpora := vars["corpora"]

	store, err := storage.NewSeparatedMemoryStorage("./data/" + corpora + ".gob")
	if err != nil {
		log.Println(err)
	}

	chatbot := &bot.ChatBot{
		LogicAdapter:   logic.NewClosestMatch(store, 1),
		Trainer:        bot.NewCorpusTrainer(store),
		StorageAdapter: store,
	}

	corporaFiles := "./data/input/trainData_" + corpora + ".json"
	_, errFile := os.Stat(corporaFiles)
	if os.IsNotExist(errFile) {
		// file does not exist, create it
		sourceFile, err := os.Open("./data/input/example.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer sourceFile.Close()

		destinationFile, err := os.Create("./data/input/trainData_temp.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destinationFile.Close()

		_, err = io.Copy(destinationFile, sourceFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.Rename("./data/input/trainData_temp.json", corporaFiles)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// write to file
	contentToWrite, err := json.Marshal(trainData)
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile(corporaFiles, contentToWrite, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// train the new data
	response.Set("success", true)
	if err := chatbot.Train(strings.Split(corporaFiles, ",")); err != nil {
		response.Set("success", false)
		response.Set("error", err.Error())
	}

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// POST none
// @tags Feed
// @Summary Endpoint to handle convert feed to training bot
// @Description Endpoint to handle convert training bot
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Res
// @Router /train/{corpora}  [post]
func ProceedFeed(w http.ResponseWriter, r *http.Request) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	response := simplejson.New()

	// Declare a new TrainData struct.
	var trainData structs.TrainData

	vars := mux.Vars(r)
	corpora := vars["corpora"]
	var result structs.Xml
	var feedUrl FeedUrl

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	errBody := json.NewDecoder(r.Body).Decode(&feedUrl)
	if errBody != nil {
		http.Error(w, errBody.Error(), http.StatusBadRequest)
		log.Printf(errBody.Error())
		return
	}

	if xmlBytes, err := getXML(feedUrl.Feed); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		xml.Unmarshal(xmlBytes, &result)
	}

	for _, item := range result.Shopitem {
		var data []string
		data = append(data, "Co mi doporučiš za produkt, prosím?")
		data = append(data, item.Product+" "+item.Imgurl+" "+item.Url)
		trainData.Conversations = append(trainData.Conversations, data)
	}

	trainData.Categories = append(trainData.Categories, "train")
	// parse to json train corpora files
	corporaFiles := "./data/input/trainData_" + corpora + ".json"
	_, errFile := os.Stat(corporaFiles)
	if os.IsNotExist(errFile) {
		// file does not exist, create it
		sourceFile, err := os.Open("./data/input/example.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer sourceFile.Close()

		destinationFile, err := os.Create("./data/input/trainData_temp.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destinationFile.Close()

		_, err = io.Copy(destinationFile, sourceFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.Rename("./data/input/trainData_temp.json", corporaFiles)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// write to file
	contentToWrite, err := json.Marshal(trainData)
	if err != nil {
		log.Printf("Marshall content: %v", err)
	}

	err = ioutil.WriteFile(corporaFiles, contentToWrite, 0644)
	if err != nil {
		log.Println(err)
	}

	// train the new data
	store, err := storage.NewSeparatedMemoryStorage("./data/" + corpora + ".gob")
	if err != nil {
		log.Printf("Failed load gob: %v", err)
	}

	chatbot := &bot.ChatBot{
		LogicAdapter:   logic.NewClosestMatch(store, 1),
		Trainer:        bot.NewCorpusTrainer(store),
		StorageAdapter: store,
	}

	response.Set("success", true)
	if err := chatbot.Train(strings.Split(corporaFiles, ",")); err != nil {
		response.Set("success", false)
		response.Set("error", err.Error())
	}

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// tweaked from: https://stackoverflow.com/a/42718113/1170664
func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}
