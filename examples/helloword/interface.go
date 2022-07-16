package main

import "fmt"

type IF interface {
	getName() string
	//getTest() string
}

type Human struct {
	firstName, lastName string
}

type Car struct {
	factory, model string
}

type Plan struct {
	core string
}

func (h *Human) getName() string {
	return h.firstName + "," + h.lastName
}

func (c *Car) getName() string {
	return c.factory + "-" + c.model
}

func (p Plan) getName() string {
	return p.core
}

func main() {
	interfaces := []IF{}
	//fmt.Printf("interfaces is %+v \n", interfaces)
	h := new(Human)
	h.firstName = "first"
	h.lastName = "last"
	interfaces = append(interfaces, h)
	c := new(Car)
	c.factory = "factory"
	c.model = "model"
	interfaces = append(interfaces, c)
	p := new(Plan)
	p.core = "A"
	//fmt.Printf("interfaces is %+v \n", interfaces)
	interfaces = append(interfaces, p)
	for _, value := range interfaces {
		fmt.Println(value.getName())
		//fmt.Println(value.getTest())
	}

}
