package fastfiber

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/nerocho/fastfiber/interf"
	"github.com/nerocho/fastfiber/utils/container"
)

var (
	configKeyPrefix = "Config_"
	configFile      = fmt.Sprintf("config.%s.yml", GetEnv("APP_ENV", "development"))
)

func initConfig() interf.ConfigInterface {
	checkFile() // 检查配置文件是否存在

	config := viper.New()
	config.AddConfigPath(".")
	config.SetConfigName(configFile) //配置文件
	config.SetConfigType("yml")      // 配置后缀

	if err := config.ReadInConfig(); err != nil {
		log.Fatal(ErrInitConfigFile + err.Error())
	}

	return &yml{config}
}

func checkFile() {
	//检查文件是否存在
	if _, err := os.Stat(configFile); err != nil {
		log.Fatal(ErrConfigFileNotExists + err.Error())
	}
}

type yml struct {
	viper *viper.Viper
}

func (y *yml) keyIsCache(keyName string) bool {
	if _, exists := container.CreateContainersFactory().KeyIsExists(configKeyPrefix + keyName); exists {
		return true
	}
	return false
}

func (y *yml) cache(keyName string, value interface{}) bool {
	return container.CreateContainersFactory().Set(configKeyPrefix+keyName, value)
}

// 通过键获取缓存的值
func (y *yml) getValueFromCache(keyName string) interface{} {
	return container.CreateContainersFactory().Get(configKeyPrefix + keyName)
}

// 清空配置项
func (y *yml) clearCache() {
	container.CreateContainersFactory().FuzzyDelete(configKeyPrefix)
}

// Get 一个原始值
func (y *yml) Get(keyName string) interface{} {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName)
	}

	value := y.viper.Get(keyName)
	y.cache(keyName, value)
	return value
}

// GetString
func (y *yml) GetString(keyName string) string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(string)
	}
	value := y.viper.GetString(keyName)
	y.cache(keyName, value)
	return value

}

// GetBool
func (y *yml) GetBool(keyName string) bool {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(bool)
	}
	value := y.viper.GetBool(keyName)
	y.cache(keyName, value)
	return value
}

// GetInt
func (y *yml) GetInt(keyName string) int {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int)
	}

	value := y.viper.GetInt(keyName)
	y.cache(keyName, value)
	return value
}

// GetInt32
func (y *yml) GetInt32(keyName string) int32 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int32)
	}

	value := y.viper.GetInt32(keyName)
	y.cache(keyName, value)
	return value
}

// GetInt64
func (y *yml) GetInt64(keyName string) int64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int64)
	} else {
		value := y.viper.GetInt64(keyName)
		y.cache(keyName, value)
		return value
	}
}

// float64
func (y *yml) GetFloat64(keyName string) float64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(float64)
	}
	value := y.viper.GetFloat64(keyName)
	y.cache(keyName, value)
	return value
}

// GetDuration
func (y *yml) GetDuration(keyName string) time.Duration {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(time.Duration)
	}

	value := y.viper.GetDuration(keyName)
	y.cache(keyName, value)
	return value
}

// GetStringSlice
func (y *yml) GetStringSlice(keyName string) []string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).([]string)
	}
	value := y.viper.GetStringSlice(keyName)
	y.cache(keyName, value)
	return value
}
