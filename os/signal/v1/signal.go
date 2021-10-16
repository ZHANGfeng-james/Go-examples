package v1

import (
	"log"
	"os"
	"os/signal"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func signalUsage() {
	log.Println("Process running...")

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt) // wait for get os.Interrupt signal

	log.Println("goroutine start to sleep")
	time.Sleep(5 * time.Second)
	log.Println("goroutine sleep over, and weak up...")

	s := <-ch
	log.Println(s)
}
