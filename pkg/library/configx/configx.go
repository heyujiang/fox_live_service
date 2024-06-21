package configx

import (
	"fmt"
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
	viper     *viper.Viper
	keyPrefix string
	cache     *configCache
}

func NewConfigX(path, prefix string, files ...string) ConfigI {
	configXOnce.Do(func() {
		config = newConfigX(path, prefix, files...)
	})

	return config
}

func newConfigX(path, prefix string, files ...string) ConfigI {
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
		viper:     configViper,
		keyPrefix: prefix,
		cache:     newConfigCache(),
	}
}

func (c *configX) Get(key string) interface{} {
	if c.cache.exist(c.getCacheConfigKey(key)) {
		return c.cache.get(c.getCacheConfigKey(key))
	} else {
		val := c.viper.Get(key)
		c.cache.set(c.getCacheConfigKey(key), val)
		return val
	}
}

func (c *configX) GetInt(key string) int {
	if c.cache.exist(c.getCacheConfigKey(key)) {
		return c.cache.get(c.getCacheConfigKey(key)).(int)
	} else {
		val := c.viper.GetInt(key)
		c.cache.set(c.getCacheConfigKey(key), val)
		return val
	}
}

func (c *configX) GetInt32(key string) int32 {
	if c.cache.exist(c.getCacheConfigKey(key)) {
		return c.cache.get(c.getCacheConfigKey(key)).(int32)
	} else {
		val := c.viper.GetInt32(key)
		c.cache.set(c.getCacheConfigKey(key), val)
		return val
	}
}

func (c *configX) GetInt64(key string) int64 {
	if c.cache.exist(c.getCacheConfigKey(key)) {
		return c.cache.get(c.getCacheConfigKey(key)).(int64)
	} else {
		val := c.viper.GetInt64(key)
		c.cache.set(c.getCacheConfigKey(key), val)
		return val
	}
}

func (c *configX) GetString(key string) string {
	if c.cache.exist(c.getCacheConfigKey(key)) {
		return c.cache.get(c.getCacheConfigKey(key)).(string)
	} else {
		val := c.viper.GetString(key)
		c.cache.set(c.getCacheConfigKey(key), val)
		return val
	}
}

func (c *configX) GetFloat64(key string) float64 {
	if c.cache.exist(c.getCacheConfigKey(key)) {
		return c.cache.get(c.getCacheConfigKey(key)).(float64)
	} else {
		val := c.viper.GetFloat64(key)
		c.cache.set(c.getCacheConfigKey(key), val)
		return val
	}
}

func (c *configX) GetIntSlice(key string) []int {
	if c.cache.exist(c.getCacheConfigKey(key)) {
		return c.cache.get(c.getCacheConfigKey(key)).([]int)
	} else {
		val := c.viper.GetIntSlice(key)
		c.cache.set(c.getCacheConfigKey(key), val)
		return val
	}
}

func (c *configX) GetStringSlice(key string) []string {
	if c.cache.exist(c.getCacheConfigKey(key)) {
		return c.cache.get(c.getCacheConfigKey(key)).([]string)
	} else {
		val := c.viper.GetStringSlice(key)
		c.cache.set(c.getCacheConfigKey(key), val)
		return val
	}
}

func (c *configX) GetStringMapString(key string) map[string]string {
	if c.cache.exist(c.getCacheConfigKey(key)) {
		return c.cache.get(c.getCacheConfigKey(key)).(map[string]string)
	} else {
		val := c.viper.GetStringMapString(key)
		c.cache.set(c.getCacheConfigKey(key), val)
		return val
	}
}

func (c *configX) GetStringMapStringSlice(key string) map[string][]string {
	if c.cache.exist(c.getCacheConfigKey(key)) {
		return c.cache.get(c.getCacheConfigKey(key)).(map[string][]string)
	} else {
		val := c.viper.GetStringMapStringSlice(key)
		c.cache.set(c.getCacheConfigKey(key), val)
		return val
	}
}

func (c *configX) GetBool(key string) bool {
	if c.cache.exist(c.getCacheConfigKey(key)) {
		return c.cache.get(c.getCacheConfigKey(key)).(bool)
	} else {
		value := c.viper.GetBool(key)
		c.cache.set(c.getCacheConfigKey(key), value)
		return value
	}
}

func (c *configX) getCacheConfigKey(key string) string {
	return fmt.Sprintf("%s_%s", c.keyPrefix, key)
}

// WatchConfig 监听文件变化
func (c *configX) WatchConfig() {
	c.viper.OnConfigChange(func(in fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if in.Op.String() == "WRITE" {
				c.cache.clear()
				lastChangeTime = time.Now()
			}
		}
	})
	c.viper.WatchConfig()
}
