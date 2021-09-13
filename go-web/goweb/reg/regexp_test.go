package reg

import "testing"

func TestIsIP(t *testing.T) {
	IsIP("12.0.0.1")
}

func TestHandleValues(t *testing.T) {
	handleValues()
}

func TestRegexp(t *testing.T) {
	regexpTest()
}

func TestRegexpExpand(t *testing.T) {
	regExpExpand()
}
