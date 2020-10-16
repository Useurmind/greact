package main

import (
	"fmt"

	"github.com/useurmind/greact/pkg/greact"
)

var switchValue = true

type RootComponent struct {

}

func (c *RootComponent) Render() greact.Element {
	switchValue, setSwitchValue := greact.UseState(true)

	renderedName := "Lise"
	if !switchValue.(bool) {
		renderedName = "Hugo"
	}

	return greact.CreateElement(
		"div", 
		greact.Props {
			"id": "greeting",
		}, 
		greact.CreateElement("div", nil, 
			greact.CreateElement("span", greact.Props{ "innerHTML": "hello span" }, nil)),
			greact.CreateElement("button", greact.Props{
				"type": "button",
				"innerHTML": "Switch Child",
				"onClick": func () {
					setSwitchValue(!switchValue.(bool))
					fmt.Println("Button pressed ", switchValue)
				},
			}),
			greact.CreateElement(&ChildComponent{
			}, &ChildComponentProps{
				greeting: renderedName,
			}, nil))
}