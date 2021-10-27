package interf

import "time"

type ConfigInterface interface {
	Get(keyName string) interface{}
	GetString(keyName string) string
	GetStringSlice(keyName string) []string
	GetInt(keyName string) int
	GetInt32(keyName string) int32
	GetInt64(keyName string) int64
	GetFloat64(keyName string) float64
	GetBool(keyName string) bool
	GetDuration(keyName string) time.Duration
}
