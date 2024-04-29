package main

import (
	"fmt"
	"runtime/metrics"
)

func main() {
	for _, v := range metrics.All() {
		fmt.Println(v.Name)
	}
}
