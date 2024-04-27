package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func main() {
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
	fmt.Printf("%+v\n", conf)
	if conf.Concurrency <= 0 {
		conf.Concurrency = 10
	}

	pool := NewPool(conf.Concurrency, conf.Duration)
	for i := 0; i < 10; i++ {
		taskId := i
		pool.Submit(func() {
			fmt.Println(taskId)
			time.Sleep(time.Second)
		})
		if i == 4 {
			time.Sleep(5 * time.Second)
		}
	}
	//time.Sleep(60 * time.Second)
	pool.Wait()
}

type Config struct {
	Concurrency int               `mapstructure:"concurrency"`
	Duration    int               `mapstructure:"duration"`
	URL         string            `mapstructure:"url"`
	Method      string            `mapstructure:"method"`
	Headers     map[string]string `mapstructure:"headers"`
}
