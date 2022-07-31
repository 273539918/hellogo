package main

import (
	"fmt"
	"golang/pkg/common"
)

func main() {

	//"[\"i-rj996a1oy510h4fwvstf\",\"123123\"]"
	arr := []string{"i-rj996a1oy510h4fwvstf", "123123"}
	r := common.StringArrayToJsonArrStr(arr)
	fmt.Println(r)

}
