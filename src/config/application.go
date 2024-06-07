package config

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/xiusin/pinecms/src/application/controllers"
	"github.com/xiusin/pinecms/src/application/models/tables"
	"github.com/xiusin/pinecms/src/common/helper"
)

const dbYml = "resources/configs/database.yml"
const appYaml = "resources/configs/application.yml"

type Config struct {
	Debug   bool     `yaml:"debug"`
	Port    int64    `yaml:"port"`
	View    viewConf `yaml:"view"`
	Session SessConf `yaml:"session"`

	LogPath      string `yaml:"log_path"`
	RuntimePath  string `yaml:"runtime_path"`
	PluginPath   string `yaml:"plugin_path"`
	PluginEnable bool   `yaml:"plugin_enable"`
	CacheDb      string `yaml:"cache_db"`

	Charset  string `yaml:"charset"`
	JwtKey   string `yaml:"jwtkey"`
	HashKey  string `yaml:"hashkey"`
	BlockKey string `yaml:"blockkey"`

	Upload struct {
		MaxBodySize int64  `yaml:"max_bodysize"`
		Engine      string `yaml:"engine"`
		BasePath    string `yaml:"base_path"`
	} `yaml:"upload"`

	Statics []struct {
		Route string `yaml:"route"`
		Path  string `yaml:"path"`
	} `yaml:"statics"`

	Search struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Url      string `yaml:"url"`
	} `yaml:"search"`
}

type SessConf struct {
	Name    string        `yaml:"name"`
	Expires time.Duration `yaml:"expires"`
}

type Html struct {
	Path   string `yaml:"path"`
	Suffix string `yaml:"suffix"`
}

type viewConf struct {
	Reload    bool   `yaml:"reload"`
	FeDirname string `yaml:"fedirname"`
	BeDirname string `yaml:"bedirname"`
	Theme     string
}

var config = &Config{}

type Site map[string]string

func (s Site) Get(key string, def ...string) string {
	val := s[key]
	if len(val) == 0 && len(def) != 0 {
		val = def[0]
	}
	return val
}

func (c *Config) init() {
	if len(c.Upload.BasePath) == 0 {
		c.Upload.BasePath = "uploads"
	}
	if len(c.RuntimePath) == 0 {
		c.RuntimePath = "runtime"
	}
	if len(c.LogPath) == 0 {
		c.LogPath = "logs"
	}
	c.LogPath = RuntimePath(c.LogPath)
	if len(c.CacheDb) == 0 {
		c.CacheDb = "cache.db"
	}
	c.CacheDb = RuntimePath(c.CacheDb)
}

func (c *Config) StaticPrefixArr() []string {
	var staticPathPrefix []string
	for _, static := range config.Statics {
		staticPathPrefix = append(staticPathPrefix, static.Route)
	}
	return staticPathPrefix
}

func IsDebug() bool {
	return config.Debug
}

func RuntimePath(path ...string) string {
	if len(path) == 0 {
		path = append(path, "")
	}
	return filepath.Join(config.RuntimePath, path[0])
}

func App() *Config {
	return config
}

func parseConfig(path string, out any) {
	absPath, err := filepath.Abs(path)
	helper.PanicErr(err)
	fileContent, err := os.ReadFile(absPath)
	helper.PanicErr(err)
	err = yaml.Unmarshal(fileContent, out)
	helper.PanicErr(err)
}

func GetSiteConfigByKey(key string, def ...string) string {
	siteConfig, _ := SiteConfig()
	return siteConfig.Get(key, def...)
}

func SiteConfig() (Site, error) {
	orm, cache := helper.GetORM(), helper.AbstractCache()
	var settingData = map[string]string{}

	err := cache.Remember(controllers.CacheSetting, &settingData, func() (any, error) {
		var settings []tables.Setting
		if err := orm.Find(&settings); err != nil {
			return nil, err
		}
		if len(settings) != 0 {
			for _, v := range settings {
				settingData[strings.ToUpper(v.Key)] = v.Value
			}
		}
		return &settingData, nil
	})
	return settingData, err
}

func init() {
	parseConfig(appYaml, config)
	config.init()
	_ = os.MkdirAll(config.LogPath, os.ModePerm)
}
