package main

import (
	"fmt"
	"reflect"
)

type T struct {
	A string
}

//注意reveive是struct，且函数名是大写
func (t *T) Ttest() string {
	return t.A + "1"
}

func main() {

	myMap := make(map[string]int, 10)
	myMap["a"] = 1
	mapT := reflect.TypeOf(myMap)
	mapV := reflect.ValueOf(myMap)
	fmt.Printf("mymap type is %+v \t , value is %+v \n", mapT, mapV)

	myStruct := T{A: "aaa"}
	structT := reflect.TypeOf(myStruct)
	valueT := reflect.ValueOf(myStruct)
	fmt.Printf("mystruct type is %+v \t,value is %+v \n", structT, valueT)

	for i := 0; i < valueT.NumField(); i++ {
		fmt.Printf("Field %d is %v \n", i, valueT.Field(i))
	}
	for i := 0; i < valueT.NumMethod(); i++ {
		fmt.Printf("Method %d is %v \n", i, valueT.Method(i))
	}

	fmt.Printf("method result is %+v \n", valueT.Method(0).Call(nil))

}
