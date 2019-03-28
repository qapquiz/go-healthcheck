package filemanager

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func GetContentFromFile(fileName string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadFile(pwd + "/" + fileName)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func IsCSVFile(fileName string) bool {
	var validCSVFile = regexp.MustCompile(`[^\s]+\.csv`)

	return validCSVFile.MatchString(fileName)
}

func ParseCSV(csvContent string) *csv.Reader {
	return csv.NewReader(strings.NewReader(csvContent))
}