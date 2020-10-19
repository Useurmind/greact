package main

import (
	"github.com/useurmind/greact/pkg/greact"
)

type ChildComponentProps struct {
	Key string
	Greeting string
}

type ChildComponent struct {
	Props *ChildComponentProps
}

func (c *ChildComponent) Render() greact.Element {
	state, setState := greact.UseState("Hello")

	return greact.CreateElement(
		"div", 
		greact.Props {
			"id": "child_comp",
		},
		greact.CreateElement("span", greact.Props{
			"key": "child_span",
			"innerHTML": state.(string) + " " + c.Props.Greeting,
		}),
		greact.CreateElement("button", greact.Props{
			"key":       "child_comp_div_button",
			"type":      "button",
			"innerHTML": "Switch Greeting",
			"onClick": func() {
				if state.(string) == "Hello" {
					setState("Goodbye")
				} else {
					setState("Hello")
				}
			},
		}),)
}
