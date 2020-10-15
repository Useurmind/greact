package main

import (
	"fmt"
	"net/http"
	"os"
)


func main() {
	addr := os.Args[1]
	assetsFolder := os.Args[2]
	fmt.Printf("Listening on %s...", addr)
	err := http.ListenAndServe(addr, http.FileServer(http.Dir(assetsFolder)))
	if err != nil {
		fmt.Println("Failed to start server", err)
		os.Exit(1)
	}
}