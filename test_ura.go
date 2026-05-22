//go:build ignore
// +build ignore

package main

import (
	"datapull/api"
	"fmt"
	"log"
)

func main() {
	rows, err := api.FetchDataURA()
	if err != nil {
		log.Fatalf("URA fetch failed: %v", err)
	}

	if len(rows) <= 1 {
		log.Fatalf("URA fetch returned no data rows")
	}

	fmt.Printf("URA rows fetched: %d\n", len(rows)-1)
	fmt.Println(rows[0])

	limit := len(rows)
	if limit > 4 {
		limit = 4
	}

	for i := 1; i < limit; i++ {
		fmt.Println(rows[i])
	}
}
