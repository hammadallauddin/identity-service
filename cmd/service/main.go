package main

import (
	"fmt"
	"log"

	"github.com/hammadallauddin/identity-service/pkg/config"
)

func main() {
	if err := config.Initialize(); err != nil {
		log.Fatalf("coult not initialize config: %v", err)
	}

	fmt.Println("Starting the service...")
}
