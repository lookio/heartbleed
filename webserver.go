package main

import (
	"github.com/hoisie/web"
)

type WebServer struct {
}

type httpHandler func(*web.Context, ...interface{})

func (webServer *WebServer) Listen() {
	web.Match("OPTIONS", "/api/v1/url", webServer.CORS)
	web.Post("/api/v1/url", webServer.checkVulnerableUrl)

	web.Run(serverConfig.Server.ListenAddress + ":" + serverConfig.Server.ListenPort)
}

func (webServer *WebServer) CORS(ctx *web.Context, args ...interface{}) {
	ctx.SetHeader("Access-Control-Allow-Origin", "*", true)
	ctx.SetHeader("Access-Control-Allow-Methods", "POST", true)
	ctx.SetHeader("Access-Control-Allow-Headers", "Access-Control-Allow-Origin", true)
	ctx.WriteString("")
}

func (webServer *WebServer) checkVulnerableUrl(ctx *web.Context, args ...interface{}) {
	sslCheck := &SslCheck{Url: ctx.Params["url"]}
	isVulnerable, err := sslCheck.CheckSync()

	ctx.SetHeader("Access-Control-Allow-Origin", "*", true)

	if err != nil {
		logger.Errorf("error checking: %v", err)
		ctx.Abort(500, "not implemented")
		return
	}

	if true == isVulnerable {
		ctx.WriteString("1")
	} else {
		ctx.WriteString("0")
	}
}
