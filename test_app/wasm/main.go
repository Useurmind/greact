package main

import (
	"fmt"
	"github.com/useurmind/greact/pkg/greact"
)

func main() {  
	fmt.Println("Go Web Assembly")
	
	greact.Render(greact.CreateElement(&RootComponent{}, nil, nil))
}

