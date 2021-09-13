package session

import (
	"container/list"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var pInstance = &ProviderInstance{
	list: list.New(),
}

func init() {
	fmt.Println("session init() in store.go")

	pInstance.sessions = make(map[string]*list.Element, 0)
	// --> manager.go Register(name string, provider Provider)
	Register("memory", pInstance)
}

func getProviderInstance() *ProviderInstance {
	return pInstance
}

type ProviderInstance struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List // a doubly linked list.
}

func (ps *ProviderInstance) SessionInit(sid string) (Session, error) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	// sid --> map[interface{}]interface{}
	v := make(map[interface{}]interface{}, 0)
	newSess := &SessionElement{
		sid:          sid,
		timeAccessed: time.Now(),
		value:        v,
	}
	element := ps.list.PushFront(newSess)
	ps.sessions[sid] = element
	return newSess, nil
}

func (ps *ProviderInstance) SessionRead(sid string) (Session, error) {
	if listElement, ok := ps.sessions[sid]; ok {
		// sessionElement 为 *list.Element 类型
		return listElement.Value.(*SessionElement), nil
	} else {
		newSess, err := ps.SessionInit(sid)
		return newSess, err
	}
}

func (ps *ProviderInstance) SessionUpdate(sid string) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	if element, ok := ps.sessions[sid]; ok {
		element.Value.(*SessionElement).timeAccessed = time.Now()
		// 更新了访问时间后，为什么 MoveToFront？ps.list 的队头节点是最新访问的 Session，队列末尾是最早访问的 Session
		ps.list.MoveToFront(element)
		return nil
	}
	return nil
}

func (ps *ProviderInstance) SessionDestory(sid string) error {
	if listElement, ok := ps.sessions[sid]; ok {
		delete(ps.sessions, sid)
		ps.list.Remove(listElement)
		return nil
	}
	return nil
}

func (ps *ProviderInstance) SessionGC(maxLifeTime int64) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	// Time 类型格式化为：2006.01.02 15:04:05
	fmt.Printf("Current: %s, maxLifeTime: %d.\n", time.Now().Format("2006.01.02 15:04:05"), maxLifeTime)

	for {
		// 为什么从 ps.list 的最后一个元素开始遍历？是最早访问的 Session
		element := ps.list.Back()
		if element == nil {
			break
		}
		// maxLifeTime 就是生命期的含义，将 Time 类型转化为了单位为秒的 int64 数值
		if (element.Value.(*SessionElement).timeAccessed.Unix() + maxLifeTime) < time.Now().Unix() {
			ps.list.Remove(element)
			delete(ps.sessions, element.Value.(*SessionElement).sid)
		} else {
			// 如果最早访问的 Session 仍在生命期，则退出遍历
			break
		}
	}
}

func (ps *ProviderInstance) print() {
	fmt.Printf("Count Session:%s.\n", strconv.Itoa(ps.list.Len()))
	index := 1
	for e := ps.list.Front(); e != nil; e = e.Next() {
		fmt.Printf("id: %d; timeAccessed: %s.\n", index, e.Value.(*SessionElement).timeAccessed)
		index++
	}
}

type SessionElement struct {
	sid          string
	timeAccessed time.Time // Session 最后访问的时间
	value        map[interface{}]interface{}
}

func (se *SessionElement) Set(key, value interface{}) error {
	se.value[key] = value
	pInstance.SessionUpdate(se.sid)
	return nil
}

func (se *SessionElement) Get(key interface{}) interface{} {
	pInstance.SessionUpdate(se.sid)
	if v, ok := se.value[key]; ok {
		return v
	} else {
		return nil
	}
}

func (se *SessionElement) Delete(key interface{}) error {
	delete(se.value, key)
	pInstance.SessionUpdate(se.sid)
	return nil
}

func (se *SessionElement) SessionID() string {
	return se.sid
}

func (se *SessionElement) String() string {
	return "sid:" + se.sid + "; timeAccessed:" + se.timeAccessed.Format("2006.01.02 15:04:05")
}
