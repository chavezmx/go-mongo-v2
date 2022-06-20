package main

import (
	"go-mongo-v2/config"
	"go-mongo-v2/routes"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"

)

func main() {
	app, err := newrelic.NewApplication(
			newrelic.ConfigAppName("go-mongo-v2"),
			newrelic.ConfigLicense("-"),
			newrelic.ConfigDistributedTracerEnabled(true),
			newrelic.ConfigDebugLogger(os.Stdout),
		)

	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(nrecho.Middleware(app))
	config.ConnectDB()
	routes.UserRoute(e)
	e.Logger.Fatal(e.Start(":6000"))
}
