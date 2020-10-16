package main

import (
	"fmt"
	"github.com/useurmind/greact/pkg/greact"
)

type ChildComponent struct {

}

func (c *ChildComponent) Render() *greact.Element {
	return greact.CreateElement(
		"div", 
		greact.Props {
			"id": "child",
			"innerHTML": "Hello child",
		})
}

type RootComponent struct {

}

func (c *RootComponent) Render() *greact.Element {
	return greact.CreateElement(
		"div", 
		greact.Props {
			"id": "greeting",
		}, 
		greact.CreateElement("div", nil, 
			greact.CreateElement("span", greact.Props{ "innerHTML": "hello span" }, nil)),
			greact.CreateElement(&ChildComponent{}, nil, nil))
}

func main() {  
	fmt.Println("Go Web Assembly")
	
	greact.Render(greact.CreateElement(&RootComponent{}, nil, nil))
}

