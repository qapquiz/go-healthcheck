package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestReadArguments_Empty(t *testing.T) {
	os.Args = []string{"go-healthcheck"}
	_, err := readArguments()

	assert.NotNil(t, err)
}

func TestReadArguments_OneArgument(t *testing.T) {
	firstArgument := "test.csv"
	os.Args = []string{"go-healthcheck", "test.csv" }
	args, err := readArguments()

	assert.Nil(t, err)

	assert.Equal(t, firstArgument, args[0])
}

func TestReadArguments_TooManyArguments(t *testing.T) {
	firstArgument := "test.csv"
	os.Args = []string{"go-healthcheck", "test.csv", "test2.csv", "test3.csv"}
	args, err := readArguments()

	assert.Nil(t, err)

	assert.Equal(t, firstArgument, args[0])
}