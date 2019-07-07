package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	timePtr := flag.Int("timer", 30, "quiz timer")
	filePtr := flag.String("file", "", "problems csv file")
	flag.Parse()

	if *filePtr == "" {
		flag.PrintDefaults()
		fmt.Println("No compatable quiz problem file found.")
		os.Exit(1)
	}

	if (*timePtr <= 0) || (*timePtr >= 200) {
		flag.PrintDefaults()
		fmt.Println("Please provide a reasonable amount of quiz time.")
		os.Exit(1)
	}

	csvData := openCSV(*filePtr)
	ch := make(chan int)
	m := make(map[string]int)

	csvToMap(m, csvData)
	var rightAnswers, wrongAnswers, index int

	go func() {
		fmt.Printf("...Press Enter to Start %d Timer...", *timePtr)
		_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
		if err != nil {
			fmt.Println("Hit ENTER to continue.")
		}
		for k, v := range m {
			badData := false
			fmt.Printf("\n Problem #%d) %s =", index, k)
			input, badData := readUserInput()
			for badData {
				fmt.Printf("Bad input, try again.")
				fmt.Printf("\n Problem #%d) %s =", index, k)
				input, badData = readUserInput()
			}
			if input == v {
				fmt.Println("\n Correct!")
				rightAnswers++
			} else {
				fmt.Println("\n Incorrect.")
				wrongAnswers++
			}
			index++
		}
		ch <- 1
	}()

	select {
	case <-ch:
		fmt.Println("\n Finished all questions in time!")
	case <-time.After(time.Duration(*timePtr) * time.Second):
		fmt.Println("\n Times Up!")
	}

	fmt.Println("\n ==== Quiz complete! ==== \n")
	fmt.Printf("Score: %.2f%% %d/%d \n", (float64(rightAnswers)/float64(len(m)))*100, rightAnswers, len(m))
}

func openCSV(filename string) *[][]string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return &lines
}

func csvToMap(questions map[string]int, csvData *[][]string) {
	for _, line := range *csvData {
		answer, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal(err)
		}
		questions[line[0]] = answer
	}
}

func readUserInput() (int, bool) {
	var input string
	_, err := fmt.Scanf("%s", &input)
	if err != nil {
		return 0, true
	} else {
		udata, err := strconv.Atoi(input)
		if err != nil {
			return 0, true
		}
		return udata, false
	}
}
