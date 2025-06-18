package confmodel

import "time"

type ServerConfig struct {
	Name string
	Port int
}

type LogConfig struct {
	Level string
}

type JWTConfig struct {
	Enable bool
	Secret string
	Expire time.Duration
}

type AppConfig struct {
	Server ServerConfig
	Log    LogConfig
	JWT    JWTConfig
}
