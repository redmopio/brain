package main

import (
	"context"
	"fmt"

	"github.com/minskylab/brain"
	"github.com/minskylab/brain/config"
)

func main() {
	config, err := config.NewLoadedConfig()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	b, err := brain.NewBrain(ctx, config)
	if err != nil {
		panic(err)
	}

	fmt.Println(b)
}
