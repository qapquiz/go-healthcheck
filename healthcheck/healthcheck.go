package healthcheck

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/qapquiz/go-healthcheck/filemanager"
)

const (
	Timeout = time.Second * 30
)

type Report struct {
	totalWebsites int
	countSuccessWebsites int
	countFailureWebsites int
}

func (report Report) print() {
	fmt.Printf("Checked websites: %d\n", report.totalWebsites)
	fmt.Printf("Successful websites: %d\n", report.countSuccessWebsites)
	fmt.Printf("Failure websites: %d\n", report.countFailureWebsites)
}

func PrintReport(report Report, totalTimeUsed float64) {
	report.print()
	fmt.Printf("Total times to finished checking website: %.4fms\n", totalTimeUsed)
}

func createClientWithTimeOut(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

func check(client *http.Client, url string, isSuccessChannel chan<- bool) {
	_, err := client.Get(url)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		isSuccessChannel <- false
	}

	// @todo error other than err.Timeout() should be handle

	if err == nil {
		isSuccessChannel <- true
	}
}

func CheckWithCSVFile(csvFileName string, sendReport chan<- Report) {
	report := Report{
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

	client := createClientWithTimeOut(Timeout)
	isSuccessChannel := make(chan bool)
	for {
		url, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err == nil {
			report.totalWebsites++
			go check(client, url[0], isSuccessChannel)
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