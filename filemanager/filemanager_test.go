package filemanager

import (
	"io"
	"testing"
)

func TestIsCSVFile(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	}{
		{"test.csv", true},
		{"test.cs", false},
		{"test", false},
	}

	for _ ,test := range tests {
		actual := IsCSVFile(test.input)
		if actual != test.expected {
			t.Errorf("Got '%v' but Expected '%v'", actual, test.expected)
		}
	}
}

func TestGetContentFromFile(t *testing.T) {
	content, err := GetContentFromFile("../test.csv")
	if err != nil {
		t.Error(err)
	}

	if len(content) == 0 {
		t.Errorf("file is empty or can't read file")
	}
}

func TestParseCSV(t *testing.T) {
	csvContent := `url
https://blognone.com
https://macthai.com
https://apple.com
`
	r := ParseCSV(csvContent)

	header, err := r.Read()
	if err != nil {
		t.Error(err)
	}

	if header[0] != "url" {
		t.Error("read error header[0] should be url")
	}

	record1, err := r.Read()
	if err != nil {
		t.Error(err)
	}

	if record1[0] != "https://blognone.com" {
		t.Error("read error record1[0] should be url")
	}

	record2, err := r.Read()
	if err != nil {
		t.Error(err)
	}

	if record2[0] != "https://macthai.com" {
		t.Error("read error record2[0] should be url")
	}

	record3, err := r.Read()
	if err != nil {
		t.Error(err)
	}

	if record3[0] != "https://apple.com" {
		t.Error("read error record3[0] should be url")
	}

	_, err = r.Read()
	if err == nil {
		t.Error("should have err io.EOF")
	}

	if err != io.EOF {
		t.Error("error when reach the end of file")
	}
}