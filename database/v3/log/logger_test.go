package log

import (
	"io/ioutil"
	"testing"
)

func TestLogger(t *testing.T) {
	SetLevel(ErrorLevel)
	if errorLog.Writer() == ioutil.Discard || infoLog.Writer() != ioutil.Discard {
		t.Fatal("errorLog's writer must be ioutil.Discard and infoLog's writer must not be ioutil.Discard")
	}

	SetLevel(Disable)
	if errorLog.Writer() != ioutil.Discard || infoLog.Writer() != ioutil.Discard {
		t.Fatal("errorLog and infoLog's writer must be ioutil.Discard")
	}
}
