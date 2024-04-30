package main

type Config struct {
	Concurrency int                 `mapstructure:"concurrency"`
	Duration    int                 `mapstructure:"duration"`
	URL         string              `mapstructure:"url"`
	Method      string              `mapstructure:"method"`
	Headers     map[string]string   `mapstructure:"headers"`
	Response    Response            `mapstructure:"response"`
	Variables   map[string]Variable `mapstructure:"variables"`
}

type Response struct {
	Code int      `mapstructure:"code"`
	Data []string `mapstructure:"data"`
}

type Variable struct {
	Source string `mapstructure:"source"`
	Length int    `mapstructure:"length"`
	Min    int    `mapstructure:"min"`
	Max    int    `mapstructure:"max"`
}
