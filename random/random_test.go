package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandStringBytesRemainder(t *testing.T) {
	randomStr1 := RandStringBytesRemainder(7)
	randomStr2 := RandStringBytesRemainder(7)

	assert.NotEqual(t, randomStr1, randomStr2)
}