package main

import (
	"errors"
	"fmt"
	"github.com/qapquiz/go-healthcheck/line"
	"os"
	"time"

	"github.com/qapquiz/go-healthcheck/filemanager"
	"github.com/qapquiz/go-healthcheck/healthcheck"
)

const lineHiringUrl = "https://hiring-challenge.appspot.com/healthcheck/report"

func readArguments() ([]string, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return nil, errors.New("there is no argument")
	}

	return args, nil
}

func main() {
	args, err := readArguments()
	if err != nil {
		fmt.Println("You must add csv file as an argument")
		fmt.Println("Example: 'go-healthcheck test.csv'")
		os.Exit(0)
	}

	csvFileName := args[0]

	if !filemanager.IsCSVFile(csvFileName) {
		fmt.Printf("'%s' is not a csv file. please try again\n", csvFileName)
		return
	}

	fmt.Println("\nPerform website checking...")

	startTime := time.Now()
	receiveReportChannel := make(chan healthcheck.Report)
	go healthcheck.CheckWithCSVFile(csvFileName, receiveReportChannel)
	healthCheckReport := <-receiveReportChannel
	totalTimeUsed := time.Since(startTime).Nanoseconds() / 1000000

	fmt.Printf("Done!\n\n")

	healthcheck.PrintReport(healthCheckReport, totalTimeUsed)

	lineAccessToken := line.GetAccessToken()
	healthcheck.SendReportToHiringLine(lineHiringUrl, lineAccessToken, healthCheckReport, totalTimeUsed)
}
