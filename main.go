package main

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	conf := parseConfig()
	pool := NewPool(conf.Concurrency, conf.Duration)
	for i := 0; i < 10; i++ {
		task := NewTask()
		task.Data = i
		pool.Submit(task)
	}
	pool.Wait()
}

// parseConfig 解析配置文件
func parseConfig() *Config {
	var configFile string
	flag.StringVar(&configFile, "c", "", "yaml文件地址，默认config.yaml")
	flag.Parse()

	if configFile == "" {
		configFile = "config.yaml"
	}
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("读取配置文件错误, %s\n", err))
	}
	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		panic(fmt.Sprintf("配置解析错误, %s\n", err))
	}
	if conf.Concurrency <= 0 {
		conf.Concurrency = 10
	}
	return &conf
}
