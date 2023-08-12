package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yiwen-ai/yiwen-static/src/app"
	"github.com/yiwen-ai/yiwen-static/src/conf"
	"github.com/yiwen-ai/yiwen-static/src/logging"
)

var help = flag.Bool("help", false, "show help info")
var version = flag.Bool("version", false, "show version info")

func main() {
	flag.Parse()
	if *help || *version {
		data, _ := json.Marshal(app.GetVersion())
		fmt.Println(string(data))
		os.Exit(0)
	}

	host := "http://" + conf.Config.Server.Addr
	logging.Infof("%s@%s start on %s", conf.AppName, conf.AppVersion, host)
	err := app.NewApp().ListenWithContext(conf.Config.GlobalSignal, conf.Config.Server.Addr)
	logging.Warningf("%s@%s http server closed: %v", conf.AppName, conf.AppVersion, err)

	ctx := conf.Config.GlobalShutdown
	for {
		if conf.Config.JobsIdle() {
			logging.Infof("%s@%s shutdown: OK", conf.AppName, conf.AppVersion)
			return
		}

		select {
		case <-ctx.Done():
			logging.Errf("%s@%s shutdown: %v", conf.AppName, conf.AppVersion, ctx.Err())
			return
		case <-time.After(time.Second):
		}
	}
}
