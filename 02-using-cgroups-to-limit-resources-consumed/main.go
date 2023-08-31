package main

import (
	"fmt"
	"strings"
	"time"
)

// the following program is to consume up to 50MB of RAM incrementally and then hold the entire program for an hour

func main() {
	var longStrs []string
	times := 50
	// build a string to consume 1MB of RAM each time it is called
	for i := 1; i <= times; i++ {
		fmt.Printf("==========%d==========\n", i)
		// let the longStrs slice hold the value/RAM consumed
		longStrs = append(longStrs, buildString(1000000, byte(i)))
	}

	// hold the application for an hour before exiting
	time.Sleep(3600 * time.Second)
}

// this buildString function builds a long string holding 'n' values
func buildString(n int, b byte) string {
	var builder strings.Builder
	// the Grow method grows to consume a specific amount of resources (in RAM)
	builder.Grow(n)
	for i := 0; i < n; i++ {
		builder.WriteByte(b)
	}
	return builder.String()
}
