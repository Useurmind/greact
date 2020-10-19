package main

import (
	"fmt"

	"github.com/useurmind/greact/pkg/greact"
)

var switchValue = true

type RootComponentProps struct {
	Key string
}

type RootComponent struct {
	Props *RootComponentProps
}

func (c *RootComponent) Render() greact.Element {
	switchValue, setSwitchValue := greact.UseState(true)
	greact.UseEffect(func() func() {
		fmt.Println("Perform Effect RootComponent", switchValue.(bool))
		return func() {
			fmt.Println("Cleanup Effect RootComponent", switchValue.(bool))
		}
	}, switchValue.(bool))

	renderedName := "Lise"
	if !switchValue.(bool) {
		renderedName = "Hugo"
	}

	return greact.CreateElement(
		"div",
		greact.Props{
			"id":  "root_comp_div",
			"key": "root_comp_div",
		},
		greact.CreateElement("div", greact.Props{"id": "root_comp_1", "key": "root_comp_div_div"},
			greact.CreateElement("span", greact.Props{"innerHTML": "hello span", "key": "root_comp_div_div_span"}, nil)),
		greact.CreateElement("button", greact.Props{
			"key":       "root_comp_div_button",
			"type":      "button",
			"innerHTML": "Switch Child",
			"onClick": func() {
				setSwitchValue(!switchValue.(bool))
				fmt.Println("Button pressed ", switchValue)
			},
		}),
		greact.CreateElement(&ChildComponent{}, &ChildComponentProps{
			Key:      "ChildComponent",
			greeting: renderedName,
		}, nil),
		func() greact.Element {
			if switchValue.(bool) {
				fmt.Println("Render conditional element")
				return greact.CreateElement("span", greact.Props{"innerHTML": "conditional span", "key": "root_comp_div_span"}, nil)
			}

			return nil
		}())
}
