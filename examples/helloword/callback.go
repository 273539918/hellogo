package main

import "fmt"

func main() {
	DoOperation(1, increase)
	DoOperation(1, decrease)
}

func DoOperation(y int, f func(int, int)) {
	f(y, 1)
}

func increase(a, b int) {
	fmt.Println("incr ease result is:", a+b)
}

func decrease(a, b int) {
	fmt.Println("decrease reseult is", a-b)
}
