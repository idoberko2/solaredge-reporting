package main

import (
	"context"

	"github.com/idoberko2/semonitor/app"
)

func main() {
	app.New().Run(context.Background())
}
