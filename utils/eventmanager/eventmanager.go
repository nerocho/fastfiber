package eventmanager

import (
	"strings"
	"sync"
)

// 事件管理中心类
type eventManage struct {
}

var sMap sync.Map

//实例化
func CreateEventManageFactory() *eventManage {
	return &eventManage{}
}

//注册事件
func (e *eventManage) Set(key string, keyFunc func(args ...interface{})) bool {
	if _, exists := e.Get(key); !exists {
		sMap.Store(key, keyFunc)
		return true
	}
	return false
}

// 获取事件
func (e *eventManage) Get(key string) (interface{}, bool) {
	if value, exists := sMap.Load(key); exists {
		return value, exists
	}
	return nil, false
}

// 调用事件
func (e *eventManage) Call(key string, args ...interface{}) {
	if valueInterface, exists := e.Get(key); exists {
		if fn, ok := valueInterface.(func(args ...interface{})); ok {
			fn(args...)
		}
	}
}

// 删除事件
func (e *eventManage) Delete(key string) {
	sMap.Delete(key)
}

// 模糊调用
func (e *eventManage) FuzzyCall(keyPre string) {
	sMap.Range(func(key, value interface{}) bool {
		if keyName, ok := key.(string); ok {
			if strings.HasPrefix(keyName, keyPre) {
				e.Call(keyName)
			}
		}
		return true
	})
}
