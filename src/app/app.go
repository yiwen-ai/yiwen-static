package app

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

	app.UseHandler(logging.AccessLogger)
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

const Wechat_UA = "MicroMessenger/"

func (gs Groups) Serve(ctx *gear.Context) error {
	logging.SetTo(ctx, "host", ctx.Host)
	if ctx.Method != http.MethodGet && ctx.Method != http.MethodHead {
		status := 200
		if ctx.Method != http.MethodOptions {
			status = 405
		}
		ctx.SetHeader(gear.HeaderContentType, "text/plain; charset=utf-8")
		ctx.SetHeader(gear.HeaderAllow, "GET, HEAD, OPTIONS")
		return ctx.End(status)
	}

	isWechat := strings.Contains(ctx.GetHeader(gear.HeaderUserAgent), Wechat_UA)
	// https://www.yiwen.pub/pub/ck1sasaglcahc6fks810?language=zho&by=ke82hfgs3ni
	if ctx.Host == "www.yiwen.pub" && !isWechat {
		next := &url.URL{
			Scheme:   "https",
			Host:     "www.yiwen.ai",
			Path:     ctx.Path,
			RawQuery: ctx.Req.URL.RawQuery,
		}
		return ctx.Redirect(next.String())
	}

	name, file := gs.lookupFile(ctx.Path)
	if name != "" {
		ctx.SetHeader(gear.HeaderCacheControl, "public, max-age=604800, must-revalidate")

		if name == "index.html" {
			ctx.SetHeader(gear.HeaderCacheControl, "no-cache, no-store")
			lang := handleContext(ctx)
			app := "web"
			if isWechat {
				app = "wechat"
			}

			html := fmt.Sprintf(`<html lang="%s" data-app="%s">`, lang, app)
			file = bytes.Replace(file, []byte("<html>"), []byte(html), 1)
		}
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
					return nil, fmt.Errorf("GetFile %q error: %v", objectKey[6:], err)
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

func handleContext(ctx *gear.Context) (lang string) {
	logging.SetTo(ctx, "referer", ctx.GetHeader(gear.HeaderReferer))
	// user preferred language
	lang = ctx.Query("lang")
	if lang == "" {
		lang = ctx.GetHeader("x-language")
	}
	if lang == "" {
		if c, _ := ctx.Req.Cookie("lang"); c != nil {
			lang = c.Value
		} else if locale := ctx.AcceptLanguage(); locale != "" {
			if i := strings.IndexAny(locale, "-_"); i > 0 {
				locale = locale[:i]
			}
			lang = locale
		}
	}

	lang = Lang639_3(lang)
	logging.SetTo(ctx, "lang", lang)

	// user preferred currency
	if cookie, _ := ctx.Req.Cookie("ccy"); cookie != nil {
		logging.SetTo(ctx, "ccy", cookie.Value)
	}

	// 用户推荐人
	if cookie, _ := ctx.Req.Cookie("by"); cookie != nil {
		logging.SetTo(ctx, "by", cookie.Value)
	} else if by := ctx.Query("by"); len(by) > 0 && len(by) <= 20 {
		// 如果 url 中包含用户推荐人，则设置到 cookie
		logging.SetTo(ctx, "by", by)
		http.SetCookie(ctx.Res, &http.Cookie{
			Name:     "by",
			Value:    by,
			HttpOnly: true,
			Secure:   conf.Config.Cookie.Secure,
			MaxAge:   int(conf.Config.Cookie.ExpiresIn),
			Path:     "/",
			Domain:   conf.Config.Cookie.Domain,
			SameSite: http.SameSiteLaxMode,
		})
	}

	return
}
