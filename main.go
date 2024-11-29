package main

import (
	"fmt"
	"github.com/ehafenmaier/boot-dev-gator/internal/config"
	"log"
)

func main() {
	// Read the configuration
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Set the user name
	err = cfg.SetUser("ehafenmaier")
	if err != nil {
		log.Fatal(err)
	}

	// Read the configuration again
	cfg, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Print the configuration to the console
	fmt.Printf("%+v\n", *cfg)
}
