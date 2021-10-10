package discover

import (
	"log"
	"reflect"
	"testing"
)

func TestDiscovery(t *testing.T) {
	servers := []string{
		"1",
		"2",
		"3",
	}

	discovery := NewMultiServersDiscovery(servers)

	all, _ := discovery.GetAll()
	if !reflect.DeepEqual(servers, all) {
		t.Fatal("discovery GetAll error")
	}

	for i := len(servers); i > 0; i-- {
		server, err := discovery.Get(RandomSelect)
		if err != nil {
			t.Fatal("discovery Get RandomSelect error")
		}
		log.Printf("Get[%d]:%s", i, server)
	}

	allCount := 30
	for i := allCount; i > 0; i-- {
		server, err := discovery.Get(RoundRobinSelect)
		if err != nil {
			t.Fatal("discovery Get RoundRobinSelect error")
		}
		log.Printf("Get[%d]:%s", i, server)

		if i == allCount/2 {
			// server 数量减少
			discovery.Update([]string{"10", "20", "30", "40", "50"})
		}
	}
}
