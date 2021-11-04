package builtin

import "testing"

func TestChannelUsage(t *testing.T) {
	channelUsage()
}

func TestChannelFromClosed(t *testing.T) {
	getEleFromClosedChannel()
}

func TestChannelForRange(t *testing.T) {
	chanForRange()
}

func TestChannelBufferRead(t *testing.T) {
	bufferChanRead()
}

func TestChannelTaskLoop(t *testing.T) {
	channelTaskLoop()
}
