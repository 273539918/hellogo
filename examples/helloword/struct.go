package main

import (
	"fmt"
	"reflect"
)

type MyType struct {
	Name    string
	TagTest string `tag1:"value1" tag2:"value2"`
}

func main() {
	t := MyType{Name: "test", TagTest: "tag"}
	printMyType(&t)

	myType := reflect.TypeOf(t)
	name := myType.Field(0)
	tagTest := myType.Field(1)
	fmt.Println(name, tagTest)
	tag := tagTest.Tag.Get("tag1")
	fmt.Println(tag)

}

func printMyType(t *MyType) {
	println(t.Name)
}
