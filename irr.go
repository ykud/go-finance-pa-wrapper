package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/alpeb/go-finance/fin"
)

// Single cashflow
type Cashflow struct {
	Date  time.Time
	Value float64
}

// Set of cashflows by key
type Cashflows map[string][]Cashflow

func createCashFlowMap(data [][]string) Cashflows {
	var cashflowMap = make(Cashflows)
	var err error
	var key string
	var date time.Time
	var value float64
	var cashflow Cashflow
	for _, line := range data {
		for j, field := range line {
			if j == 0 {
				key = field
			} else if j == 1 {
				date, err = time.Parse(time.DateOnly, field)
				if err != nil {
					fmt.Printf("Error converting date: %v", err)
				}
				cashflow.Date = date
			} else if j == 2 {
				value, err = strconv.ParseFloat(field, 64)
				if err != nil {
					fmt.Printf("Error converting value from string: %v", err)
				}
				cashflow.Value = value
			}
		}
		cashflowMap[key] = append(cashflowMap[key], cashflow)
	}
	return cashflowMap
}
func main() {
	var inputFile = flag.String("in", "", "Input file with cashflows in the key,cashflow_value csv format")
	var outputFile = flag.String("out", "", "Output file with XIRRs in the key,XIRR csv format")
	var functionToCalculate = flag.String("function", "irr", "Function to execute IRR / XIRR")
	var guess = flag.Float64("guess", 0, "Guess parameter for function")
	var logFileName = flag.String("log", "irr.log", "Log file")

	flag.Parse()
	logFile, err := os.OpenFile(*logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	defer logFile.Close()
	log.Printf(
		`Starting %v calculation with input file %v and output file %v, guess %v`,
		*functionToCalculate,
		*inputFile,
		*outputFile,
		*guess,
	)
	start := time.Now()
	inFile, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	csvReader := csv.NewReader(inFile)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	cashFlowMap := createCashFlowMap(data)
	outFile, err := os.OpenFile(
		*outputFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	csvWriter := csv.NewWriter(outFile)

	var csvIRRs [][]string
	var result float64
	for key, value := range cashFlowMap {
		var cfs []float64
		var dates []time.Time
		result = 0
		for _, val := range value {
			cfs = append(cfs, val.Value)
			dates = append(dates, val.Date)
		}

		if *functionToCalculate == "xirr" {
			result, err = fin.ScheduledInternalRateOfReturn(cfs, dates, *guess)
			if err != nil {
				log.Printf("%v %v", key, err)
			}
		} else if *functionToCalculate == "irr" {
			result, err = fin.InternalRateOfReturn(cfs, *guess)
			if err != nil {
				log.Printf("%v %v", key, err)
			}
		} else {
			log.Printf("Unknown function %s", *functionToCalculate)
		}
		row := []string{key, strconv.FormatFloat(result, 'f', 8, 64)}
		csvIRRs = append(csvIRRs, row)
	}
	csvWriter.WriteAll(csvIRRs)
	elapsed := time.Since(start)
	log.Printf(
		"%v calculation took %s, calculated %d IRRs",
		*functionToCalculate,
		elapsed,
		len(cashFlowMap),
	)
}
