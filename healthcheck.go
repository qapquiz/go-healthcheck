package main

import "fmt"

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

}

func checkWebsiteInCSVFile(csvFileName string, sendReport chan<- HealthCheckReport) {
	report := HealthCheckReport{
		totalWebsites: 0,
		countSuccessWebsites: 0,
		countFailureWebsites: 0,
	}

	var url string
	isSuccessChannel := make(chan bool)
	// read csv file somehow
	for i := 0; i < 10; i++ {
		report.totalWebsites++
		go healthCheck(url, isSuccessChannel)
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