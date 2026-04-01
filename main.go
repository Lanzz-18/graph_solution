package main

import (
	"fmt"
	"os"
)

func main(){
	if len(os.Args) < 2 {
		fmt.Println()
		os.Exit(1)
	}

	filePath := os.Args[1]
	fmt.Println("Reading file:", filePath)
}