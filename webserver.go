package main

import (
	"github.com/hoisie/web"
)

type WebServer struct {
}

type httpHandler func(*web.Context, ...interface{})

func (webServer *WebServer) Listen() {
        web.Post("/api/v1/url", webServer.checkVulnerableUrl)

	web.Run(serverConfig.Server.ListenAddress + ":" + serverConfig.Server.ListenPort)
}

func (webServer *WebServer) hasAuth(f httpHandler) httpHandler {
	return func(ctx *web.Context, data ...interface{}) {
		f(ctx, nil)
	}
}

func (webServer *WebServer) checkVulnerableUrl(ctx *web.Context, args ...interface{}) {
        sslCheck := &SslCheck{Url: ctx.Params["url"]}
        isVulnerable, err := sslCheck.CheckSync()

        if err != nil {
            logger.Errorf("error checking: %v", err)
            ctx.Abort(500, "not implemented")
            return
        }

        if (true == isVulnerable) {
            ctx.WriteString("1")
        } else {
            ctx.WriteString("0")
        }
}
