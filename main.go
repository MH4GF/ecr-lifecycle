package main

import (
	"fmt"
	"os"
)

// build時に注入する
var Revision string

func main() {
	app := newApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
