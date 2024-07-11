package main

import (
	"fmt"
	"os"
)

func main() {
	err := Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", rootCmd.Use, err.Error())
		os.Exit(1)
	}
}
