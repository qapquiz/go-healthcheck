package filemanager

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCSVFile(t *testing.T) {
	tests := map[string]struct{
		input string
		expected bool
	}{
		"test.csv must be valid": {input: "test.csv", expected: true},
		"test.cs must be invalid": {input: "test.cs", expected: false},
		"test must be invalid": {input: "test", expected: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := IsCSVFile(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestGetContentFromFile(t *testing.T) {
	pathToFile := "../test.csv"
	content, err := GetContentFromFile(pathToFile)

	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(content))
}

func TestGetContentFromFile_MustFailed_When_FileNotFound(t *testing.T) {
	_, err := GetContentFromFile("../not_found_file.csv")
	assert.NotNil(t, err)
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