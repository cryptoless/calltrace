package router

import (
	"calltrace/app/api"

	"calltrace/config"
	"fmt"
	"time"

	"github.com/cryptoless/chain-raw-api-server/message"
	"github.com/cryptoless/chain-raw-api-server/util/response"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"golang.org/x/time/rate"
)

func RouteInit() {
	limit := rate.Every(time.Duration(config.RateCfg.Interval) * time.Millisecond)
	rateLimit = rate.NewLimiter(limit, config.RateCfg.Burst)

	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddleRateLimit)
		group.ALL("/call", api.TraceApi.Call)
	})

	s.BindStatusHandlerByMap(map[int]ghttp.HandlerFunc{
		403: func(r *ghttp.Request) { r.Response.Write("403, status", r.Get("status"), " found") },
		404: func(r *ghttp.Request) { r.Response.Write("404, status", r.Get("status"), " found") },
		500: func(r *ghttp.Request) { r.Response.Write("500, status", r.Get("status"), " found") },
	})
}

var rateLimit *rate.Limiter

func MiddleRateLimit(r *ghttp.Request) {
	if !rateLimit.Allow() {
		err := message.InternalError(fmt.Sprintf("rateLimit unAllow"))
		g.Log().Error(err)
		response.ErrorResponse(r, err)
	}
	r.Middleware.Next()
}
