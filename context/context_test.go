package main

import (
	"testing"
)

func TestContext(t *testing.T) {
	withValueUseNormalType()
}

func TestCancelContext(t *testing.T) {
	cancelContextPropagate()
}
