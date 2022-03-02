package main

import (
	_ "calltrace/boot"
	"calltrace/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Log().SetAsync(true)
	router.RouteInit()
	g.Server().Run()
}
