package sync

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/go-examples-with-tests/tools"
)

func conditionUsage() {
	c := sync.NewCond(&sync.Mutex{})
	var ready int

	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Int63n(5)) * time.Second)

			// 加锁更改等待条件
			c.L.Lock()
			ready++
			c.L.Unlock()

			log.Printf("运动员#%d 已准备就绪\n", i)
			// 广播唤醒所有的等待者
			c.Broadcast()
		}(i)
	}

	c.L.Lock()
	for ready != 10 {
		c.Wait()
		log.Println("裁判员被唤醒一次")
	}
	c.L.Unlock()

	//所有的运动员是否就绪
	log.Println("所有运动员都准备就绪。比赛开始，3，2，1, ......")
}

type Queue struct {
	cond *sync.Cond
	cap  int
	data []interface{}
}

func NewQueue(cap int) *Queue {
	return &Queue{
		cond: &sync.Cond{L: &sync.Mutex{}},
		data: make([]interface{}, 0),
		cap:  cap,
	}
}

func (q *Queue) Enqueue(ele interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for len(q.data) == q.cap {
		q.cond.Wait() // 必须要先 q.cond.L.Lock()
	}

	q.data = append(q.data, ele)
	q.cond.Broadcast()
}

func (q *Queue) Dequeue() interface{} {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for len(q.data) == 0 {
		q.cond.Wait() // 必须要先 q.cond.L.Lock()
	}

	ele := q.data[0]
	// 移除对头的元素
	q.data = q.data[1:]
	log.Printf("len:%d, cap:%d", len(q.data), cap(q.data))
	q.cond.Broadcast()
	return ele
}

func (queue *Queue) Info() {
	log.Println(tools.SliceInfo("Queue", queue.data))
}

func QueueUsage() {
	queue := NewQueue(3)

	go func() {
		time.Sleep(1 * time.Second)
		ele := queue.Dequeue()
		log.Printf("dequeue element:%v", ele)
	}()

	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)
	queue.Info()

	queue.Enqueue(4)

	ele := queue.Dequeue()
	log.Printf("dequeue ele:%v", ele)
}
