package main

import (
	"context"

	"github.com/robertobadjio/tgtime-router-tracker/internal/app"
	"github.com/robertobadjio/tgtime-router-tracker/internal/logger"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatal("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		logger.Fatal("failed to run app: %s", err.Error())
	}
}
