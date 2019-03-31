package healthcheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net"
	"net/http"
	"os"
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

func PrintReport(report Report, totalTimeUsed int64) {
	report.print()
	fmt.Printf("Total times to finished checking website: %d ms\n", totalTimeUsed)
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
		return
	}

	if err != nil {
		isSuccessChannel <- false
		return
	}

	isSuccessChannel <- true
}

func CheckWithCSVFile(csvFileName string, sendReport chan<- Report) {
	report := Report{
		totalWebsites: 0,
		countSuccessWebsites: 0,
		countFailureWebsites: 0,
	}

	csvContent, err := filemanager.GetContentFromFile(csvFileName)
	if err != nil {
		fmt.Printf("%s. please try again\n", err)
		os.Exit(1)
	}

	csvReader := filemanager.ParseCSV(csvContent)
	_, err = csvReader.Read() // skip header
	if err == io.EOF {
		fmt.Printf("'%s' is empty. please try again\n", csvFileName)
		os.Exit(1)
	}


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

func SendReportToHiringLine(url string, lineAccessToken *oauth2.Token, report Report, totalTimeUsed int64) {
	client := createClientWithTimeOut(Timeout)
	body := map[string]int64{
		"total_websites": int64(report.totalWebsites),
		"success": int64(report.countSuccessWebsites),
		"failure": int64(report.countFailureWebsites),
		"total_time": totalTimeUsed,
	}

	jsonValue, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {

	}
	req.Header.Add("Authorization", "Bearer " + lineAccessToken.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	fmt.Println("sending report...")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("cannot send report to %s. please try again\n", url)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("get error status: %d from %s. please try again\n" ,resp.StatusCode, url)
		os.Exit(1)
	}

	fmt.Println("send report successfully!")
}