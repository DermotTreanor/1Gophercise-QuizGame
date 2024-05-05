package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"errors"
)

var filename *string = flag.String("csv", "problem.csv", "this is the help message for this flag")
var total_questions, total_right, total_wrong int 

func main(){
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil{
		fmt.Printf("The following error occured: %v\n", err)
		return
	}	
	
	//Since os.File is a type that implements io.Reader we can pass it here.
	var csv_reader *csv.Reader = csv.NewReader(file)
	var csv_line []string
	for{
		csv_line, err = csv_reader.Read()
		if err == io.EOF{
			fmt.Println("That's the end of the quiz!")
			break
		}else if err != nil{
			fmt.Fprintf(os.Stderr, "There is an UNEXPECTED error: %v\n", err)
			continue
		}
		//Check that we can get an answer from the question field. Skip the line if we can't. 
		result, err := solve_question(csv_line[0])
		if err != nil{
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		//Increment the valid question total and prompt the user with the whole question field. 
		total_questions++
		fmt.Println(csv_line[0])
		var given_input string
		fmt.Scanf("%s", &given_input)
		mark_attempt(result, given_input)
	}
	//Display the results of the questions. 
	fmt.Println(total_questions, total_right, total_wrong)
}


func solve_question(question string) (int, error){
	//Find the math portion of the question with a regex
	reg := regexp.MustCompile(`\d+[\+\-\*/]\d+`)
	finding := reg.Find([]byte(question))
	if len(finding) < 1{
		return 0, errors.New("Couldn't find appropriate math expression.")
	}
	
	//Use a new regex to get the operator and operands. Convert the operands to ints.
	var operand_ints []int
	operator_regex := regexp.MustCompile("[*-+/]")
	operator := string(operator_regex.Find(finding))
	operand_strings := operator_regex.Split(string(finding), -1)

	for _, v := range operand_strings{
		op, err := strconv.Atoi(v)
		if err != nil{
			return 0, nil
		}
		operand_ints = append(operand_ints, op)
	}

	//Perform the correct math calculation and return the result.
	switch operator{
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


func mark_attempt(true_answer int, attempt string){
	int_attempt, err := strconv.Atoi(strings.Trim(attempt, " "))
	if err != nil{
		total_wrong++
		return
	}
	if true_answer == int_attempt{
		total_right++
	}else{
		total_wrong++
	}
	return
}

