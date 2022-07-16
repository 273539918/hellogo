package main

import (
	"encoding/json"
	"fmt"
)

type D struct {
	FirstName string
	LastName  string
}

func (d D) API_test() string {
	return d.FirstName + "1"
}

func main() {

	t := D{"first", "last"}
	//struct to json string
	tStr := Struct2JsonString(t)
	fmt.Printf("t str is %+v \n", tStr)
	// json to struct
	tStruct := JonString2Struct(tStr)
	fmt.Printf("t struct is %+v \n", tStruct)

	var obj interface{}
	err := json.Unmarshal([]byte(tStr), &obj)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("t struct is %+v \n", obj)
	}

	objMap, _ := obj.(map[string]interface{})
	for k, v := range objMap {
		fmt.Println(k, v)
		switch value := v.(type) {
		case string:
			fmt.Printf("type of %s is string,value is %v \n", k, value)
		case interface{}:
			fmt.Printf("type of %s is interface{}, value is %v \n", k, value)
		default:
			fmt.Printf("type of %s is wrong,value is %v \n", k, value)
		}
	}
}

func Struct2JsonString(t D) string {
	strByte, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
	}
	return string(strByte)
}

func JonString2Struct(jsonStr string) D {
	d := D{}
	err := json.Unmarshal([]byte(jsonStr), &d)
	if err != nil {
		fmt.Println(err)
	}
	return d
}
