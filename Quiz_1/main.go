package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Problem struct that stores each question
type Problem struct {
	question string
	answer   string
}

//Function to register the score
func registerScore(p Problem, ans chan int, i int) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("P #%d: %s => ", i, p.question)
	scanner.Scan()
	text := scanner.Text()
	text = strings.ReplaceAll(text, " ", "")
	if text == p.answer {
		ans <- 1
	} else {
		ans <- 0
	}
}

func main() {

	// Flags for the program
	filename := flag.String("file", "problems.csv", "name of the csv file containing problems")
	timerTime := flag.Int("time", 3, "alloted time to complete the Quiz")
	shuffle := flag.Bool("shuffle", false, "shuffle the problems, false by default")
	flag.Parse()

	//Opening and Reading CSV file
	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer file.Close()

	f := csv.NewReader(file)
	var questions []Problem

	for {
		row, err := f.Read()
		prob := Problem{}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		prob.question = row[0]
		prob.answer = strings.ReplaceAll(row[1], " ", "")
		questions = append(questions, prob)
	}

	totalQuestions := len(questions)

	//Exit if there are no problems in the CSV file
	if totalQuestions == 0 {
		fmt.Println("Warning : chosen csv file is empty!\nExiting...")
		time.Sleep(time.Second * 1)
		os.Exit(0)
	}

	//Shuffle the problems if flag is true
	rand.Seed(time.Now().Unix())
	if *shuffle {
		rand.Shuffle(totalQuestions, func(i, j int) {
			questions[i], questions[j] = questions[j], questions[i]
		})
	}

	score := 0

	//create channel for the answers
	ans := make(chan int)
	t := *timerTime

	//creates a go routine and a channel and excecutes for given seconds
	timer := time.NewTimer(time.Second * time.Duration(t))
	defer timer.Stop()

	//Iterate through all the questions
	for i, p := range questions {

		//go routine to register the answer from user
		go registerScore(p, ans, i+1)
		select {
		case res := <-ans:
			if res == 1 {
				score++
			}
		// if the timer get a value it concludes that alloted time is finished and program exits
		case <-timer.C:
			fmt.Println("\nTime Out !!! \nFinal Score :", score,
				"\nYou got", score, "correct answer(s) out of", totalQuestions, "questions.")
			os.Exit(0)
		}
	}

	fmt.Println("Final Score :", score, "\nYou got", score, "correct answer(s) out of", totalQuestions, "questions.")
}
