package config

import (
	"os"
	"strconv"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
)

var RateCfg rateCfg

var CmdParser *gcmd.Parser

func CfgInit() {
	CmdParser, err := gcmd.Parse(g.MapStrBool{
		"c,config": true,
	})
	if err != nil {
		panic(err)
	}
	if err := g.Cfg().SetPath(CmdParser.GetOpt("c", ".")); err != nil {
		panic(err)
	}

	port := os.Getenv("port")
	if port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			panic(err)
		}
		g.Server().SetPort(p)
	}

	(&RateCfg).Load()

}
