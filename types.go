package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Config struct {
	Concurrency int               `mapstructure:"concurrency"`
	Duration    int               `mapstructure:"duration"`
	URL         string            `mapstructure:"url"`
	Method      string            `mapstructure:"method"`
	Headers     map[string]string `mapstructure:"headers"`
	Response    Response          `mapstructure:"response"`
}

type Response struct {
	Code int      `mapstructure:"code"`
	Data []string `mapstructure:"data"`
}

type Task struct {
	id     int
	config Config
}

func NewTask(id int, conf Config) *Task {
	return &Task{
		id:     id,
		config: conf,
	}
}

func (t *Task) Exec() {
	client := &http.Client{}
	req, err := http.NewRequest(strings.ToUpper(t.config.Method), t.config.URL, bytes.NewReader([]byte{}))
	if err != nil {
		fmt.Printf("http.NewRequest failed, error: %s", err)
		return
	}
	for k, v := range t.config.Headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("client.Do failed, error: %s", err)
		return
	}
	if resp.StatusCode != t.config.Response.Code {
		fmt.Println("status code 不一致")
		return
	}
	defer resp.Body.Close()
	byt, _ := io.ReadAll(resp.Body)
	data := string(byt)
	total := len(t.config.Response.Data)
	matchNum := 0
	for _, str := range t.config.Response.Data {
		if strings.Contains(data, str) {
			matchNum += 1
		}
	}
	if total == matchNum {
		fmt.Printf("%d 执行结果: %d\n", t.id, resp.StatusCode)
	} else {
		fmt.Println("不匹配")
	}
}
