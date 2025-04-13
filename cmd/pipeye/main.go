package main

import (
	"fmt"
	"os"

	"github.com/cpaluszek/pipeye/internal/app"
	"github.com/cpaluszek/pipeye/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Config: ", cfg)
	}

	pipeyeApp := app.New(cfg)
	if err := pipeyeApp.Run(); err != nil {
		fmt.Printf("Error running app: %v\n", err)
		os.Exit(1)
	}
}
