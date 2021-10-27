package container

import (
	"strings"
	"sync"
)

var sMap sync.Map

type containers struct{}

//容器工厂
func CreateContainersFactory() *containers {
	return &containers{}
}

func (c *containers) Set(k string, v interface{}) (res bool) {
	if _, exists := c.KeyIsExists(k); !exists {
		sMap.Store(k, v)
		res = true
	}
	return
}

func (c *containers) Get(k string) interface{} {
	if v, exists := c.KeyIsExists(k); exists {
		return v
	}
	return nil
}

func (c *containers) Delete(k string) {
	sMap.Delete(k)
}

// 按照key前缀模糊删除
func (c *containers) FuzzyDelete(keyPrefix string) {
	sMap.Range((func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			if strings.HasPrefix(k, keyPrefix) {
				sMap.Delete(k)
			}
		}
		return true
	}))
}

func (c *containers) KeyIsExists(k string) (interface{}, bool) {
	return sMap.Load(k)
}
