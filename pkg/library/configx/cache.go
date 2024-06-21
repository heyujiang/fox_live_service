package configx

import "sync"

type configCache struct {
	sync.Map
}

func newConfigCache() *configCache {
	return &configCache{}
}

func (cc *configCache) set(key string, value interface{}) {
	cc.Store(key, value)
}

func (cc *configCache) get(key string) interface{} {
	v, _ := cc.Load(key)
	return v
}

func (cc *configCache) delete(key string) {
	cc.Delete(key)
}

func (cc *configCache) exist(key string) bool {
	_, ok := cc.Load(key)
	return ok
}

func (cc *configCache) clear() {
	cc.Range(func(key, value interface{}) bool {
		cc.delete(key.(string))
		return true
	})
}
