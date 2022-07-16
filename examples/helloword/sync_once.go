package main

import (
	"fmt"
	"sync"
)

type SliceNum []int

func NewSlice() SliceNum {
	return make(SliceNum, 0)
}

func (s *SliceNum) Add(item int) *SliceNum {

	*s = append(*s, item)
	fmt.Println("add", item)
	fmt.Println("add SliceNum end", *s)
	return s
}

func main() {
	once := sync.Once{}
	s := NewSlice()
	//多个Do的调用只会执行一次
	once.Do(func() {
		s.Add(16)
	})
	once.Do(func() {
		s.Add(16)
	})
	once.Do(func() {
		s.Add(16)
	})
}
