package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	acc, err := NewAccount("a", "b", "password")
	fmt.Printf("%+v\n", acc)

	assert.Nil(t, err)
}
