package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"
)

func readArguments() ([]string, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return nil, errors.New("there is no argument")
	}

	return args, nil
}

func isCSVFile(fileName string) bool {
	var validCSVFile = regexp.MustCompile(`[^\s]+\.csv`)

	return validCSVFile.MatchString(fileName)
}

func main() {
	args, err := readArguments()
	if err != nil {
		fmt.Println("You must add csv file as an argument")
		fmt.Println("Example: 'go-healthcheck test.csv'")
		os.Exit(0)
	}

	csvFileName := args[0]

	if !isCSVFile(csvFileName) {
		fmt.Printf("'%s' is not a csv file. please try again\n", csvFileName)
		return
	}

	fmt.Println("Perform website checking...")

	startTime := time.Now()
	receiveReportChannel := make(chan HealthCheckReport)
	checkWebsiteInCSVFile(csvFileName, receiveReportChannel)
	healthCheckReport := <- receiveReportChannel
	totalTimeUsed := time.Since(startTime).Seconds() * 1000

	fmt.Println("Done!")

	printReport(healthCheckReport, totalTimeUsed)

	//lineAccessToken := getLineAccessToken()
	//sendReport(sendReportUrl, lineAccessToken)
}
