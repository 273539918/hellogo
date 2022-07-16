package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	//unsafeWrite()
	safeWrite()
	time.Sleep(time.Second * 1)
}

//多次运行会报错，fatal error: concurrent map writes
func unsafeWrite() {
	configMap := map[string]int{}
	for i := 0; i < 100; i++ {
		go func() {
			configMap["A"] = i
		}()
	}
	fmt.Println(configMap["A"])
}

type SafeMap struct {
	safeMap map[string]int
	sync.Mutex
}

func (s *SafeMap) Write(k string, v int) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.safeMap[k] = v
}

func safeWrite() {
	s := SafeMap{
		safeMap: map[string]int{},
		Mutex:   sync.Mutex{},
	}
	fmt.Println(s)
	for i := 0; i < 100; i++ {
		go func() {
			s.Write("A", i)
		}()
	}
	fmt.Println(s.safeMap["A"])
}
