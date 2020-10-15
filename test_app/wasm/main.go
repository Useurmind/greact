package main

import (
	"fmt"
	"syscall/js"
)

func main() {  
	fmt.Println("Go Web Assembly")
	
	jsDoc := js.Global().Get("document")
	if !jsDoc.Truthy() {
		handleError(fmt.Errorf("Could not retrieve js document"))
		return
	}

	jsBody := jsDoc.Get("body")
	if !jsBody.Truthy() {
		handleError(fmt.Errorf("Could not retrieve js body"))
		return
	}

	jsRoot := jsDoc.Call("createElement", "div")
	if !jsRoot.Truthy() {
		handleError(fmt.Errorf("Could not create js root"))
		return
	}
	jsRoot.Set("id", "root")
	jsRoot.Set("innerHTML", "Hello WASM")

	jsBody.Call("appendChild", jsRoot)

}

func handleError(err error) {
	fmt.Errorf("ERROR: %v", err)
}