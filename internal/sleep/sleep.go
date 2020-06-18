package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func printUsage() {
	fmt.Printf("Usage: %s <sleep_seconds>", filepath.Base(os.Args[0]))
}

func main() {
	if len(os.Args) != 2 {
		printUsage()
		return
	}
	sleepDurSec, err := strconv.Atoi(os.Args[1])
	if err != nil {
		printUsage()
		os.Exit(1)
	}
	time.Sleep(time.Duration(sleepDurSec) * time.Second)
}
