package main

import (
	"fmt"
	"os"
)

func main() {
	app := newApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
