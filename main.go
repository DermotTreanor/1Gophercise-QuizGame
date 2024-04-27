package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

var filename *string = flag.String("csv", "problem.csv", "this is the help message for this flag")

func main(){
	flag.Parse()
	file, err := os.Open(*filename)
	if err != nil{
		fmt.Printf("The following error occured: %v\n", err)
	}	
	
	//Since os.File is a type that implements io.Reader we can pass it here.
	//This basically means it is a type that the NewReader can read data from to make its own Reader.
	var csv_reader *csv.Reader = csv.NewReader(file)
	var csv_line []string
	for true{
		csv_line, err = csv_reader.Read()
		if err == io.EOF{
			fmt.Println("That's the end of the quiz!")
			break
		}else if err != nil{
			fmt.Printf("There is an UNEXPECTED error: %v\n", err)
		}
		fmt.Println("What is " + csv_line[0]+ ": ")
		var answer string = "%s"
		fmt.Scanf(answer)
		fmt.Println(answer)
	}
	
}