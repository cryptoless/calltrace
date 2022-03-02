package api

import (
	"calltrace/app/service"

	"github.com/cryptoless/chain-raw-api-server/message"
	"github.com/cryptoless/chain-raw-api-server/util/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var TraceApi = traceApi{}

type traceApi struct{}

// Index is a demonstration route handler for output "Hello World!".
func (*traceApi) Call(r *ghttp.Request) {
	_, err := r.WebSocket()
	if err != nil {
		// http
		msg, err := message.ParseMessage(r.GetBody())
		if err != nil {
			g.Log().Error(err)
			response.ErrorResponse(r, err)
		}
		g.Log().Debug("Api:", msg.Method)

		//
		rst := service.Service.HandleMsg(r.GetCtx(), msg)
		response.Response(r, rst)
	} else {
		// ws, until to close
		// service.WsCon(r.GetCtx(), ws).Poll()
		g.Log().Debug("no ws.")
	}
}
