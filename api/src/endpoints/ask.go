package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	bot "main/bot"
	"main/bot/adapters/logic"
	structs "main/structs/bot"
	utils "main/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

// POST none
// @tags Ask
// @Summary Endpoint to handle communication with bot
// @Description Endpoint to handle communication with bot
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ResAsk
// @Router /ask  [post]
func ProceedAsk(w http.ResponseWriter, r *http.Request, storageStruct []structs.StorageForAI) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	response := simplejson.New()

	// Declare a new Question struct.
	var question structs.Question

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	errBody := json.NewDecoder(r.Body).Decode(&question)
	if errBody != nil {
		http.Error(w, errBody.Error(), http.StatusBadRequest)
		log.Printf(errBody.Error())
		return
	}

	// set chatbot instance
	if question.Corpora == "" {
		question.Corpora = "corpora"
	}

	var store structs.StorageForAI
	for _, m := range storageStruct {
		if m.Name == question.Corpora {
			store = m
		}
	}

	chatbot := &bot.ChatBot{
		LogicAdapter:   logic.NewClosestMatch(store.Storage, 1),
		Trainer:        bot.NewCorpusTrainer(store.Storage),
		StorageAdapter: store.Storage,
	}

	startTime := time.Now()

	if question.IsTeach {
		corporaFiles := "./data/input/trainData_" + question.Corpora + ".json"
		_, err := os.Stat(corporaFiles)
		if os.IsNotExist(err) {
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

		// save the quesition to file
		// If the file doesn't exist, create it, or append to the file
		f, err := os.OpenFile("no_answer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write([]byte("--" + question.Text + "\n")); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}

		// load json to learn
		content, err := ioutil.ReadFile(corporaFiles)
		if err != nil {
			log.Fatal(err)
		}

		row := structs.TrainData{}
		err = json.Unmarshal(content, &row)

		// update
		add := []string{question.Question, question.Text}
		row.Conversations = append(row.Conversations, add)

		// write to file
		contentToWrite, err := json.Marshal(row)
		if err != nil {
			fmt.Println(err)
		}
		err = ioutil.WriteFile(corporaFiles, contentToWrite, 0644)
		if err != nil {
			log.Fatal(err)
		}

		// train the new data
		if err := chatbot.Train(strings.Split(corporaFiles, ",")); err != nil {
			log.Fatal(err)
		}

		response.Set("answer", "Thank you bro, I have learned something new :)")
		response.Set("success", true)
		response.Set("needLearn", false)

	} else {
		answers := chatbot.GetResponse(question.Text)
		if len(answers) == 0 {
			response.Set("answer", "I am not able to answer your question :( Please teach me how to answer on '"+question.Text+"'.")
			response.Set("success", true)
			response.Set("needLearn", true)
			// save the quesition to file
			// If the file doesn't exist, create it, or append to the file
			f, err := os.OpenFile("no_answer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			if _, err := f.Write([]byte(question.Text + "\n")); err != nil {
				log.Fatal(err)
			}
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		} else {
			response.Set("answer", answers[0].Content)
			response.Set("success", true)
			response.Set("needLearn", false)
		}
	}
	response.Set("answerTimeMicroS", time.Since(startTime).Microseconds())

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
