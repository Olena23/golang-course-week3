package main

import (
	"time"
	"github.com/Sirupsen/logrus"
	"github.com/tcnksm/go-input"
	"os"
	"log"
	"strconv"
	"fmt"
	"encoding/csv"
	"bufio"
	"io"
	"strings"
)

var question = ""
var response = ""
var newResponse = ""
var goOn = true
var cities []string

func checkLastLetter (s string, cities []string) {
	for j :=1; j < 6 ; j++ {
		if (len(s)-j>=0) && ((len(s)-j+1)>=0) {
			var lastLetter = s[len(s)-j:len(s)-j+1]
			for i := 0; i < len(cities); i++ {
				if cities[i][:1] == lastLetter {
					newResponse = cities[i]
					cities = append(cities[:i], cities[(i+1):]...)
					if (response != newResponse){
						response = newResponse
						return
					}
				}
			}
		}
	}
	goOn = false
}

func citiesGame (pinger chan string, id int, name string, cities []string) {
	question = <-pinger
	checkLastLetter (question, cities)
	pinger <- response

	logrus.WithFields(logrus.Fields{
		"id":       id,
		"name":     name,
		"question": question,
		"response": response,
	}).Info("New Round")
	}

func main() {
	//Prompt for number of players
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	query := "How many players do we need? Type number from 2 to 13"
	quantity, err := ui.Ask(query, &input.Options{
		Default: "2",
		Required: true,
		Loop:     true,
		ValidateFunc: func(s string) error {
			number, _ := strconv.Atoi(s)
			if (number >= 2) && (number <= 13) {
				return nil
			}
			return fmt.Errorf("input must be from 2 to 13")
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	//number of players
	players, err := strconv.Atoi(quantity)

	playerNames := []string {"Anna", "Max", "Ivy", "Paul", "Suzy", "Oleg", "Angela", "Tom", "Olga", "Svetlana", "Peter", "Dymitriy", "Val"}

	//Form list of city names
	csvFile, _ := os.Open("cities.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		addName := strings.Trim(line[1], " ")
		cities = append(cities, addName)
	}

	//delete header from cities
	cities = append(cities[:0], cities[1:]...)

	gameField := make(chan string, 1)
	gameField <- cities[0]
	for {
			for i := 1; i <= players; i++ {
				if goOn {
					go citiesGame (gameField, i, playerNames[i-1], cities)
					time.Sleep(time.Second/200)
				} else {
					logrus.WithFields(logrus.Fields{
						"Game status":       "Game Over",
					}).Info("This game is over")
					return
				}
			}
	}

	// The main goroutine starts by sending into the channel
	for {
		// Block the main thread until an interrupt
		time.Sleep(time.Second)
	}
}