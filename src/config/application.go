package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/xiusin/pine"
	"github.com/xiusin/pine/cache"
	"github.com/xiusin/pinecms/src/application/controllers"
	"github.com/xiusin/pinecms/src/application/models/tables"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug   bool    `yaml:"debug"`
	Favicon string  `yaml:"favicon"`
	Port    int64   `yaml:"port"`
	View    View    `yaml:"view"`
	Session Session `yaml:"session"`
	LogPath string  `yaml:"log_path"`
	CacheDb string  `yaml:"cache_db"`
	Statics []struct {
		Route string `yaml:"route"`
		Path  string `yaml:"path"`
	} `yaml:"statics"`
	Charset           string `yaml:"charset"`
	JwtKey            string `yaml:"jwtkey"`
	HashKey           string `yaml:"hashkey"`
	BlockKey          string `yaml:"blockkey"`
	BackendRouteParty string `yaml:"backend_route_party"`
	Upload            struct {
		MaxBodySize int64  `yaml:"max_bodysize"`
		Engine      string `yaml:"engine"`
		BasePath    string `yaml:"base_path"`
	} `yaml:"upload"`
}

type Session struct {
	Name    string        `yaml:"name"`
	Expires time.Duration `yaml:"expires"`
}

type Html struct {
	Path   string `yaml:"path"`
	Suffix string `yaml:"suffix"`
}

type View struct {
	Reload    bool   `yaml:"reload"`
	FeDirname string `yaml:"fedirname"`
	BeDirname string `yaml:"bedirname"`
	Theme     string
}

const dbYml = "resources/configs/database.yml"
const appYml = "resources/configs/application.yml"

var config = &Config{}

func init() {
	parseConfig(appYml, config)
	_ = os.MkdirAll(config.LogPath, os.ModePerm)
}

func AppConfig() *Config {
	return config
}

func parseConfig(path string, out interface{}) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	fileContent, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err.Error())
	}
	err = yaml.Unmarshal(fileContent, out)
	if err != nil {
		panic(err.Error())
	}
}

func SiteConfig() (map[string]string, error) {
	xorm := pine.Make(controllers.ServiceXorm).(*xorm.Engine)
	cache := pine.Make(controllers.ServiceICache).(cache.AbstractCache)
	var settingData = map[string]string{}
	err := cache.GetWithUnmarshal(controllers.CacheSetting, &settingData)
	if err != nil {
		var settings []tables.Setting
		err := xorm.Find(&settings)
		if err != nil {
			return nil, err
		}
		if len(settings) != 0 {
			for _, v := range settings {
				settingData[v.Key] = v.Value
			}
		}
		if err = cache.SetWithMarshal(controllers.CacheSetting, &settingData); err != nil {
			return nil, err
		}
	}
	return settingData, nil
}
