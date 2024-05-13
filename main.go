package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	//"io"
	"bufio"
	"errors"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var filename *string = flag.String("csv", "problem.csv", "this is the help message for this flag")
var clock *int = flag.Int("clock", 30, "selects the number of seconds the quiz lasts")
var shuffle *bool = flag.Bool("shfl", false, "shuffle the order of the quiz used")
var total_questions, total_right, total_wrong int

var mu sync.Mutex

func main() {
	flag.Parse()
	go run_quiz()

	time.Sleep(time.Duration(*clock) * time.Second)
	end_quiz()
}

func run_quiz() {
	//Open the file specified with the flag
	file, err := os.Open(*filename)
	if err != nil {
		fmt.Printf("The following error occured: %v\n", err)
		return
	}

	//Since os.File is a type that implements io.Reader we can pass it here.
	var csv_reader *csv.Reader = csv.NewReader(file)
	var csv_lines [][]string

	//Read all the records to a slice. Shuffle the slice if flag is set
	csv_lines, err = csv_reader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in CSV file: \n%v\n", err)
		fmt.Fprintln(os.Stderr, "Exiting...")
		os.Exit(1)
	}
	if len(csv_lines) < 1 {
		fmt.Fprint(os.Stderr, "Error: No records found in file.")
	}
	if *shuffle {
		shuffle_slice(csv_lines)
	}

	//Loop through the slice to ask each question.
	for _, csv_line := range csv_lines {
		//Check that we can get an answer from the question field. Skip the line if we can't.
		result, err := solve_question(csv_line[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		//Increment the valid question total and prompt the user with the whole question field.
		mu.Lock()
		total_questions++
		mu.Unlock()
		fmt.Println(csv_line[0])
		var given_input string

		//Query: Why is it that we cannot use the fmt.Scan functions here without the users extra spaces/values/newlines causing issues.
		input := bufio.NewScanner(os.Stdout)
		input.Scan()
		given_input = input.Text()
		mu.Lock()
		mark_attempt(result, given_input)
		mu.Unlock()

	}
	end_quiz()
}

func shuffle_slice(sl [][]string) {

	jumbled := make([][]string, len(sl))

	//Create the source of random numbers and the Rand object that will use them.
	seed := time.Now().Unix()
	source := rand.NewSource(seed)
	rand_point := rand.New(source)

	//Create a slice of jumbled inds and use it to put sl values in jumbled randomly
	mixed_inds := rand_point.Perm(len(sl))
	for i, v := range mixed_inds {
		jumbled[i] = sl[v]
	}
	copy(sl, jumbled)
}

func solve_question(question string) (int, error) {
	//Find the math portion of the question with a regex
	reg := regexp.MustCompile(`\d+[\+\-\*/]\d+`)
	finding := reg.Find([]byte(question))
	if finding == nil {
		return 0, errors.New("Couldn't find appropriate math expression.")
	}

	//Use a new regex to get the operator and operands. Convert the operands to ints.
	var operand_ints []int
	operator_regex := regexp.MustCompile("[*-+/]")
	operator := string(operator_regex.Find(finding))
	operand_strings := operator_regex.Split(string(finding), -1)

	for _, v := range operand_strings {
		op, err := strconv.Atoi(v)
		if err != nil {
			return 0, nil
		}
		operand_ints = append(operand_ints, op)
	}

	//Perform the correct math calculation and return the result.
	switch operator {
	case "*":
		return operand_ints[0] * operand_ints[1], nil
	case "/":
		return operand_ints[0] / operand_ints[1], nil
	case "+":
		return operand_ints[0] + operand_ints[1], nil
	case "-":
		return operand_ints[0] - operand_ints[1], nil
	default:
		return 0, errors.New("Unexpected parsing of question field.")
	}
}

func mark_attempt(true_answer int, attempt string) {
	//Trim the spacing and convert to an int. If we encounter an error when converting then we increment wrong total
	int_attempt, err := strconv.Atoi(strings.Trim(attempt, " "))
	if err != nil {
		//fmt.Println("WRONG INPUT!")
		total_wrong++
		return
	}
	if true_answer == int_attempt {
		//fmt.Println("RIGHT Answer!")
		total_right++
	} else {
		//fmt.Println("WRONG Answer!")
		total_wrong++
	}
	return
}

func end_quiz() {
	mu.Lock()
	if total_questions > (total_wrong + total_right) {
		total_questions--
	}
	fmt.Println(total_questions, total_right, total_wrong) //Display results
	mu.Unlock()
	os.Exit(0)
}
