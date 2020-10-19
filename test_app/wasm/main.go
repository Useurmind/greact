package main

import (
	"fmt"

	"github.com/useurmind/greact/pkg/greact"
	"github.com/useurmind/greact/pkg/greactdom"
)

func main() {  
	fmt.Println("Go Web Assembly")
	
	greactdom.Render(greact.CreateElement(&RootComponent{}, &RootComponentProps{ Key: "RootComponent" }, nil))
}

