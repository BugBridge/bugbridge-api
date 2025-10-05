package main

import (
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/BugBridge/bugbridge-api/api/handlers"
	"github.com/BugBridge/bugbridge-api/config"
)

func main() {
	a := handlers.App{}
	a.Config = *config.New()

	err := a.Initialize() //initialize database and router
	if err != nil {
		zap.S().With(err).Error("error calling initialize")
		return
	}

	// Add database middleware
	routerWithDB := handlers.DatabaseMiddleware(a.GetDBHelper())(a.Router)

	zap.S().Infow("BugBridge API is up and running", "url", a.Config.BaseURL, "port", a.Config.Port)
	log.Fatal(http.ListenAndServe(":"+a.Config.Port, routerWithDB))
}
