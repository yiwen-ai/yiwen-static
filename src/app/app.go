package app

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/teambition/gear"

	"github.com/yiwen-ai/yiwen-static/src/conf"
	"github.com/yiwen-ai/yiwen-static/src/logging"
	"github.com/yiwen-ai/yiwen-static/src/service"
)

var startTime = time.Now()

// NewApp ...
func NewApp() *gear.App {
	app := gear.New()

	app.Set(gear.SetTrustedProxy, true)
	// ignore TLS handshake error
	app.Set(gear.SetLogger, log.New(gear.DefaultFilterWriter(), "", 0))
	app.Set(gear.SetCompress, gear.ThresholdCompress(128))
	app.Set(gear.SetGraceTimeout, time.Duration(conf.Config.Server.GracefulShutdown)*time.Second)
	app.Set(gear.SetEnv, conf.Config.Env)

	groups, err := LoadFiles(conf.Config.GlobalSignal)
	if err != nil {
		logging.Panicf("Load files error: %v", err)
	}

	app.UseHandler(groups)
	return app
}

type Group struct {
	Prefix  string
	Default []byte
	Files   map[string][]byte
}

type Groups []*Group

func (gs Groups) lookupFile(path string) (string, []byte) {
	for _, group := range gs {
		if strings.HasPrefix(path, group.Prefix) {
			name := path[len(group.Prefix):]
			if data, ok := group.Files[name]; ok {
				return name, data
			}

			if len(group.Default) > 0 {
				return "index.html", group.Default
			}
		}
	}

	return "", nil
}

func (gs Groups) Serve(ctx *gear.Context) error {
	if ctx.Method != http.MethodGet && ctx.Method != http.MethodHead {
		status := 200
		if ctx.Method != http.MethodOptions {
			status = 405
		}
		ctx.SetHeader(gear.HeaderContentType, "text/plain; charset=utf-8")
		ctx.SetHeader(gear.HeaderAllow, "GET, HEAD, OPTIONS")
		return ctx.End(status)
	}

	name, file := gs.lookupFile(ctx.Path)
	if name != "" {
		http.ServeContent(ctx.Res, ctx.Req, name, startTime, bytes.NewReader(file))
	}

	return nil
}

func LoadFiles(ctx context.Context) (Groups, error) {
	oss := service.NewOSS()

	var groups []*Group
	for prefix, files := range conf.Config.Assets {
		if prefix == "" || len(files) == 0 {
			return nil, fmt.Errorf("invalid assets config: %#v", conf.Config.Assets)
		}

		group := &Group{
			Prefix: prefix,
			Files:  make(map[string][]byte),
		}

		for name, objectKey := range files {
			if strings.HasPrefix(objectKey, "oss://") {
				data, err := oss.GetFile(ctx, objectKey[6:])
				if err != nil {
					return nil, err
				}

				logging.Infof("Load %s from: %s, %d bytes", prefix+name, objectKey, len(data))

				if name == "*" {
					group.Default = data
				} else {
					group.Files[name] = data
				}
			} else {
				return nil, fmt.Errorf("unsupported objectKey: %s", objectKey)
			}
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func GetVersion() map[string]string {
	return map[string]string{
		"name":      conf.AppName,
		"version":   conf.AppVersion,
		"buildTime": conf.BuildTime,
		"gitSHA1":   conf.GitSHA1,
	}
}
