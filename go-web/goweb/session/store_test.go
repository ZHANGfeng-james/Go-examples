package session

import (
	"fmt"
	"testing"
	"time"
)

func TestSessionInint(t *testing.T) {
	globalSession, _ := NewManager("memory", "gosessionid", 3600)
	pInstance := getProviderInstance()
	if globalSession == nil || pInstance == nil {
		t.Fatal("Provider Instance is nil!")
	}

	MAX_COUNT := 3
	for index := 0; index < MAX_COUNT; index++ {
		sid := globalSession.sessionId()
		session, err := pInstance.SessionInit(sid)
		if err != nil {
			t.Fatal("Create Session error!")
		}
		fmt.Printf("sid: %s; session: %s.\n", sid, session)

		if index == MAX_COUNT-1 {
			break
		}
		time.Sleep(5 * time.Second)
	}

	pInstance.print()

	pInstance.SessionGC(6)

	pInstance.print()

}
