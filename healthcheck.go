package main

import (
	"fmt"
	"github.com/qapquiz/go-healthcheck/filemanager"
	"io"
)

type HealthCheckReport struct {
	totalWebsites int
	countSuccessWebsites int
	countFailureWebsites int
}

func (report HealthCheckReport) print() {
	fmt.Printf("Checked websites: %d\n", report.totalWebsites)
	fmt.Printf("Successful websites: %d\n", report.countSuccessWebsites)
	fmt.Printf("Failure websites: %d\n", report.countFailureWebsites)
}

func printReport(report HealthCheckReport, totalTimeUsed float64) {
	report.print()
	fmt.Printf("Total times to finished checking website: %.2fms", totalTimeUsed)
}

func healthCheck(url string, isSuccessChannel chan<- bool) {
	//MakeGetRequest(url)
}

func checkWebsiteInCSVFile(csvFileName string, sendReport chan<- HealthCheckReport) {
	report := HealthCheckReport{
		totalWebsites: 0,
		countSuccessWebsites: 0,
		countFailureWebsites: 0,
	}

	csvContent, err := filemanager.GetContentFromFile(csvFileName)
	if err != nil {
		fmt.Printf("reading '%s' error. please try again\n", err)
	}

	csvReader := filemanager.ParseCSV(csvContent)
	csvReader.Read() // skip header

	isSuccessChannel := make(chan bool)
	for {
		url, err := csvReader.Read()
		if err == nil {
			go healthCheck(url[0], isSuccessChannel)
		}

		if err == io.EOF {
			break
		}
	}

	for i := 0; i < report.totalWebsites; i++ {
		isSuccess := <-isSuccessChannel
		if isSuccess {
			report.countSuccessWebsites++
		} else {
			report.countFailureWebsites++
		}
	}

	sendReport <- report
}