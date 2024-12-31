package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s seconds [exit-code]\n", os.Args[0])
		os.Exit(2)
	}

	seconds, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Println("Invalid number of seconds:", os.Args[1])
		os.Exit(2)
	}

	exitCode := 0
	if len(os.Args) > 2 {
		exitCode, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid exit code:", os.Args[2])
			os.Exit(exitCode)
		}
	}

	time.Sleep(time.Duration(seconds * float64(time.Second)))
	fmt.Println("T")
	os.Exit(exitCode)
}
