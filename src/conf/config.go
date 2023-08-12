package conf

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/teambition/gear"
)

// Config ...
var Config ConfigTpl

var AppName = "yiwen-static"
var AppVersion = "0.1.0"
var BuildTime = "unknown"
var GitSHA1 = "unknown"

var once sync.Once

func init() {
	p := &Config
	readConfig(p, "../../config/default.toml")
	if err := p.Validate(); err != nil {
		panic(err)
	}
	p.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	p.GlobalSignal = gear.ContextWithSignal(context.Background())

	var cancel context.CancelFunc
	p.GlobalShutdown, cancel = context.WithCancel(context.Background())
	go func() {
		<-p.GlobalSignal.Done()
		time.AfterFunc(time.Duration(p.Server.GracefulShutdown)*time.Second, cancel)
	}()
}

type Logger struct {
	Level string `json:"level" toml:"level"`
}

type Server struct {
	Addr             string `json:"addr" toml:"addr"`
	GracefulShutdown uint   `json:"graceful_shutdown" toml:"graceful_shutdown"`
}

type OSS struct {
	Bucket          string `json:"bucket" toml:"bucket"`
	Endpoint        string `json:"endpoint" toml:"endpoint"`
	AccessKeyId     string `json:"access_key_id" toml:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret" toml:"access_key_secret"`
}

// ConfigTpl ...
type ConfigTpl struct {
	Rand           *rand.Rand
	GlobalSignal   context.Context
	GlobalShutdown context.Context
	Env            string                       `json:"env" toml:"env"`
	Logger         Logger                       `json:"log" toml:"log"`
	Server         Server                       `json:"server" toml:"server"`
	OSS            OSS                          `json:"oss" toml:"oss"`
	Assets         map[string]map[string]string `json:"assets" toml:"assets"`

	globalJobs int64 // global async jobs counter for graceful shutdown
}

func (c *ConfigTpl) Validate() error {
	return nil
}

func (c *ConfigTpl) ObtainJob() {
	atomic.AddInt64(&c.globalJobs, 1)
}

func (c *ConfigTpl) ReleaseJob() {
	atomic.AddInt64(&c.globalJobs, -1)
}

func (c *ConfigTpl) JobsIdle() bool {
	return atomic.LoadInt64(&c.globalJobs) <= 0
}

func readConfig(v interface{}, path ...string) {
	once.Do(func() {
		filePath, err := getConfigFilePath(path...)
		if err != nil {
			panic(err)
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		_, err = toml.Decode(string(data), v)
		if err != nil {
			panic(err)
		}
	})
}

func getConfigFilePath(path ...string) (string, error) {
	// 优先使用的环境变量
	filePath := os.Getenv("CONFIG_FILE_PATH")

	// 或使用指定的路径
	if filePath == "" && len(path) > 0 {
		filePath = path[0]
	}

	if filePath == "" {
		return "", fmt.Errorf("config file not specified")
	}

	return filePath, nil
}