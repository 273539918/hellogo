package main

import (
	"fmt"
	"sync"
)

func main() {

	var sm sync.Map
	// store 方法，添加元素
	sm.Store(1, "a")
	// load 方法，获取value
	if v, ok := sm.Load(1); ok {
		fmt.Println(v)
	}
	//LoadOrStore方法，如果该key存在且没有被标记删除，则返回value和true
	//如果该key不存在，则存储新的value，并返回value和false
	v1, ok1 := sm.LoadOrStore(1,
		"b")
	fmt.Printf("value is %s , result is %v \n", v1, ok1)
	v2, ok2 := sm.LoadOrStore(2, "c")
	fmt.Printf("value is %s , result is %v \n", v2, ok2)

	//遍历该map，参数是个函数，该函数的两个参数是遍历获得的key和value，返回值是一个bool值，返回false表示遍历立即结束
	sm.Range(func(k, v interface{}) bool {
		fmt.Printf("%v:%v\n", k, v)
		return true
	})

}
