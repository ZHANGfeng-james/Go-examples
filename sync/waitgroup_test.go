package sync

import "testing"

func TestWaitGroup(t *testing.T) {
	waitGroup()
}

func TestWaitGroupReuse(t *testing.T) {
	waitGroupReuse()
}

func TestWaitGroupGetCount(t *testing.T) {
	calcWaitGroupCount()
}
