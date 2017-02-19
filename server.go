package webhookproxy

import (
	log "github.com/Sirupsen/logrus"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/grokify/webhook-proxy-go/src/adapters"
	"github.com/grokify/webhook-proxy-go/src/config"
	"github.com/grokify/webhook-proxy-go/src/handlers"
	"github.com/grokify/webhook-proxy-go/src/handlers/slack"
	"github.com/grokify/webhook-proxy-go/src/handlers/travisci"
)

const (
	ROUTE_SLACK_IN_GLIP           = "/webhook/slack/in/glip/:webhookuid"
	ROUTE_SLACK_IN_GLIP_SLASH     = "/webhook/slack/in/glip/:webhookuid/"
	ROUTE_TRAVISCI_OUT_GLIP       = "/webhook/travisci/out/glip/:webhookuid"
	ROUTE_TRAVISCI_OUT_GLIP_SLASH = "/webhook/travisci/out/glip/:webhookuid/"
)

// StartServer initializes and starts the webhook proxy server
func StartServer(cfg config.Configuration) {
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(cfg.LogLevel)

	adapter, err := adapters.NewGlipAdapter("")
	if err != nil {
		panic("Cannot build adapter")
	}

	router := fasthttprouter.New()

	router.GET("/", handlers.HomeHandler)

	slackInHandler := slack.NewSlackToGlipHandler(cfg, &adapter)
	router.POST(ROUTE_SLACK_IN_GLIP, slackInHandler.HandleFastHTTP)
	router.POST(ROUTE_SLACK_IN_GLIP_SLASH, slackInHandler.HandleFastHTTP)

	travisciOutHandler := travisci.NewTravisciOutToGlipHandler(cfg, &adapter)
	router.POST(ROUTE_TRAVISCI_OUT_GLIP, travisciOutHandler.HandleFastHTTP)
	router.POST(ROUTE_TRAVISCI_OUT_GLIP_SLASH, travisciOutHandler.HandleFastHTTP)

	log.Fatal(fasthttp.ListenAndServe(cfg.Address(), router.Handler))
}
