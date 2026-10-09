package main

import (
	"fmt"
	"log"

	config "trash_deleter/internal/config"
	scanner "trash_deleter/internal/scanner"
)

func main() {
	cfg := config.Parse()

	items, err := scanner.Scan(cfg.Root, cfg.Patterns, cfg.IngoreDirs)
	if err != nil {
		log.Fatal(err)
	}

	for _, it := range items {
		fmt.Println(it)
	}
}
