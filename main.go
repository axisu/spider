package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789*_-"
const stringLength = 3

func generateString(l int) string {
	s := make([]string, 0, l)
	attempts := 0
	for i := 0; i < l; i++ {
		attempts++
		randomIndex := rand.Intn(len(charset))
		s = append(s, string(charset[randomIndex]))
	}
	return strings.Join(s, "")
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			target := "279fxp*zMTwCpLohBTy_1fu_JPAiSmskuh*rj88ktbYCHRaGNFimH13siXf7icB5TSIVVtfWlRiAdcG2a_nj8A0"
			l := len(target)
			attempts := 0
			for {
				attempts++
				generated := generateString(l)
				if generated == target {
					break
				}
			}
			fmt.Printf("生成目标字符串所需尝试次数：%d\n", attempts)
			return
		}()
	}
	wg.Wait()
	return

	conf := parseConfig()
	//fmt.Printf("%+v\n", conf)

	pool := NewPool(conf.Concurrency, conf.Duration)
	for i := 0; i < 10; i++ {
		task := NewTask(i, conf)
		pool.Submit(task)
	}
	pool.Wait()
}

// parseConfig 解析配置文件
func parseConfig() Config {
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
	return conf
}

func parseConfig1() Config {
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
	return conf
}
