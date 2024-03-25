package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/thatisuday/commando"

	"github.com/idoberko2/semonitor/app"
)

const ArgDays = "days"

func main() {
	commando.
		SetExecutableName("semonitor").
		SetVersion("v1.0.0").
		SetDescription("This tool fetches SolarEdge energy statistics and stores them to timescale db")

	a := app.New()
	commando.
		Register(nil).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			if err := a.RunServer(context.Background()); err != nil {
				log.WithError(err).Fatal("error running app server")
			}
		})

	commando.
		Register("fetch-recent").
		AddArgument(ArgDays, "how many days to fetch (starting today going backwards)", "").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			if err := a.RunLastDays(context.Background(), args[ArgDays].Value); err != nil {
				log.WithError(err).Fatal("error running app fetch recent")
			}
		})

	commando.Parse(nil)
}
