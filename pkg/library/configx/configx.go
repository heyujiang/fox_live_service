package configx

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"
)

var lastChangeTime time.Time
var config ConfigI
var configXOnce sync.Once

func init() {
	lastChangeTime = time.Now()
}

type configX struct {
	viper *viper.Viper
}

func NewConfigX(path string, files ...string) ConfigI {
	configXOnce.Do(func() {
		config = newConfigX(path, files...)
	})

	return config
}

func newConfigX(path string, files ...string) ConfigI {
	configViper := viper.New()
	configViper.AddConfigPath(path)
	if len(files) > 0 {
		configViper.SetConfigName(files[0])
	} else {
		configViper.SetConfigName("config")
	}
	configViper.SetConfigType("yaml")

	if err := configViper.ReadInConfig(); err != nil {
		log.Fatalf("config file read error:%s", err)
	}
	return &configX{
		viper: configViper,
	}
}

func (c *configX) Get(key string) interface{} {
	return c.viper.Get(key)
}

func (c *configX) GetInt(key string) int {
	return c.viper.GetInt(key)
}

func (c *configX) GetInt32(key string) int32 {
	return c.viper.GetInt32(key)
}

func (c *configX) GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

func (c *configX) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *configX) GetFloat64(key string) float64 {
	return c.viper.GetFloat64(key)
}

func (c *configX) GetIntSlice(key string) []int {
	return c.viper.GetIntSlice(key)
}

func (c *configX) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

func (c *configX) GetStringMapString(key string) map[string]string {
	return c.viper.GetStringMapString(key)
}

func (c *configX) GetStringMapStringSlice(key string) map[string][]string {
	return c.viper.GetStringMapStringSlice(key)
}

func (c *configX) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// WatchConfig 监听文件变化
func (c *configX) WatchConfig() {
	c.viper.OnConfigChange(func(in fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if in.Op.String() == "WRITE" {
				lastChangeTime = time.Now()
			}
		}
	})
	c.viper.WatchConfig()
}
