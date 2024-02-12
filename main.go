package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"

	"webpush/config"
	"webpush/transport"
	"webpush/util"
)

func main() {

	configs := config.LoadConfig(config.LoadEnvConfig())

	e := echo.New()
	e.Debug = true // full log

	subscriptionManager := util.NewJSONSubscriptionManager()
	pageHandler := transport.NewPageHandler(configs)
	webPushHandler := transport.NewWebPushHandler(configs, subscriptionManager)

	e.Static("/", "assets")
	e.GET(transport.LANDING_PAGE, pageHandler.GetLandingPage)
	e.GET(transport.GET_VAPID_PUBLIC_KEY, webPushHandler.GetVapidPubKey)
	e.POST(transport.SAVE_SUBCRIPTION, webPushHandler.SaveSubscription)
	e.POST(transport.SEND_NOTIFICATION, webPushHandler.SendNotification)

	e.Use(transport.BuildMiddlewareLogger())

	port := configs.GetENVConfigs().ServerPort
	address := fmt.Sprintf(":%v", port)
	slog.Info("Listening", slog.String("port", port))
	e.Logger.Fatal(e.Start(address))
}
