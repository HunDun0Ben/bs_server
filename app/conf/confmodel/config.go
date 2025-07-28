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
	Enable        bool
	Secret        string
	Expire        time.Duration
	RefreshExpire time.Duration
}

type AppConfig struct {
	Server ServerConfig
	Log    LogConfig
	JWT    JWTConfig
	Redis  RedisConfig
}

// RedisConfig 定义了 Redis 的配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}
