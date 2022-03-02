package config

import "github.com/gogf/gf/frame/g"

type rateCfg struct {
	Interval int
	Burst    int
}

func (a *rateCfg) Validate() error {

	return nil
}
func (a *rateCfg) Load() {
	a.Interval = g.Cfg().GetInt("server.Interval")
	a.Burst = g.Cfg().GetInt("server.Burst")
}
