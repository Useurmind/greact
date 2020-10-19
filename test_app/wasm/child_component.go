package main

import (
	"github.com/useurmind/greact/pkg/greact"
)

type ChildComponentProps struct {
	Key string
	greeting string
}

type ChildComponent struct {
	Props *ChildComponentProps
}

func (c *ChildComponent) Render() greact.Element {
	return greact.CreateElement(
		"div", 
		greact.Props {
			"id": "child_comp",
			"innerHTML": "Hello " + c.Props.greeting,
		})
}
