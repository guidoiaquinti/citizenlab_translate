package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckIfString(t *testing.T) {
	var input interface{}
	var result bool

	input = "string"
	result = checkIfString(input)
	assert.True(t, result)

	input = map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	result = checkIfString(input)
	assert.False(t, result)
}
