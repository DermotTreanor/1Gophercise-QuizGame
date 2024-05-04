package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

func get_math_expression(question string) string{
	reg, err := regexp.Compile("[0-9]+[\\+\\-\\*/][0-9]+")
	if err != nil{
		fmt.Println(err)
	}
	findings := reg.Find([]byte(question))
	return string(findings)
}

var filename *string = flag.String("csv", "problem.csv", "this is the help message for this flag")

func main(){
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil{
		fmt.Printf("The following error occured: %v\n", err)
	}	
	
	//Since os.File is a type that implements io.Reader we can pass it here.
	var csv_reader *csv.Reader = csv.NewReader(file)
	var csv_line []string
	var total_questions, right_ans, wrong_ans int 
	for{
		csv_line, err = csv_reader.Read()
		if err == io.EOF{
			fmt.Println("That's the end of the quiz!")
			break
		}else if err != nil{
			fmt.Printf("There is an UNEXPECTED error: %v\n", err)
			continue
		}
		total_questions++
		fmt.Println(get_math_expression(csv_line[0]))
		
		var given_input string
		fmt.Scanf("%s", &given_input)
		fmt.Println(given_input)
	}

	fmt.Println(total_questions, right_ans, wrong_ans)
	
}