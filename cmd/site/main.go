package main

import (
	"log"

	"cricidev/site/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
