package main

import "fmt"

func main() {
	myMap := make(map[string]string, 10)
	myMap["a"] = "b"
	myMap["c"] = "d"
	myFuncMap := map[string]func() int{
		"funcA": func() int { return 1 },
	}
	fmt.Println(myMap)
	fmt.Println(myFuncMap)
	f := myFuncMap["funcA"]
	fmt.Println(f())
	value, exists := myMap["a"]
	if exists {
		fmt.Println(value)
	}
	for k, v := range myMap {
		fmt.Println(k, v)
	}

}
