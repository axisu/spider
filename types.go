package main

import (
	"fmt"
	"time"
)

type Config struct {
	Concurrency int               `mapstructure:"concurrency"`
	Duration    int               `mapstructure:"duration"`
	URL         string            `mapstructure:"url"`
	Method      string            `mapstructure:"method"`
	Headers     map[string]string `mapstructure:"headers"`
}

type Task struct {
	Data any
}

func NewTask() *Task {
	return &Task{}
}

func (t *Task) Exec() {
	fmt.Println(t.Data)
	time.Sleep(time.Second)
}
