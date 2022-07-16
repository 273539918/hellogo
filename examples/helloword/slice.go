package main

import "fmt"

func main() {

	myArray := [5]int{0, 1, 2, 3, 4}
	for _, value := range myArray {
		value *= 2
	}
	//0，1，2，3，4
	fmt.Printf("myArray %+v\n", myArray)
	for index, _ := range myArray {
		myArray[index] *= 2
	}
	// 0，2，4，6，8
	fmt.Printf("myArray %+v\n", myArray)
	mySlice := myArray[1:3]
	//0，2，4，6，8
	fmt.Printf("myArray %+v\n", myArray)
	//2，4
	fmt.Printf("mySlice %+v\n", mySlice)
	mySlice = append(mySlice, 5)
	//0，2，4，5，8
	fmt.Printf("myArray %+v\n", myArray)
	//2，4，5
	fmt.Printf("mySlice %+v\n", mySlice)
	fullSlice := myArray[:]
	fmt.Printf("address of fullSlice %p \n", fullSlice)
	fmt.Printf("address of myArray %p \n", &myArray)
	//0，2，4，5，8
	fmt.Printf("fullSlice %+v\n", fullSlice)
	remove3rdItem := deleteItem(fullSlice, 1)
	//0,4,5,8
	fmt.Printf("remove3rdItem %+v\n", remove3rdItem)
	//0,4,5,8,8
	fmt.Printf("myArray %+v\n", myArray)
	fmt.Printf("fullSlice %+v\n", fullSlice)
	sliceAddItem := addItem(fullSlice)
	fmt.Printf("sliceAddItem %+v\n", sliceAddItem)
	fmt.Printf("myArray %+v\n", myArray)
	fmt.Printf("fullSlice %+v\n", fullSlice)

}

func deleteItem(slice []int, index int) []int {
	fmt.Printf("%+v\n", slice)
	//append方法，如果修改后的大小不操作slice,会直接将结果也修改到第一个参数中
	tmp := append(slice[:index], slice[index+1:]...)
	fmt.Printf("%+v\n", slice)
	return tmp
}

func addItem(slice []int) []int {
	fmt.Printf("%+v\n", slice)
	//append方法，会直接将结果也修改到第一个参数中
	tmp := append(slice[:], 0, 0, 0)
	fmt.Printf("%+v\n", slice)
	return tmp
}
