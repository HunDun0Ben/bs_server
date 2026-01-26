package confmodel

import "time"

type ServerConfig struct {
	Name string // 服务器名称
	Port int    // 服务器端口
}

type LogConfig struct {
	Level string // 日志级别
}

type JWTConfig struct {
	Enable        bool          // 是否启用 JWT
	Secret        string        // 密钥
	Expire        time.Duration // 过期时间
	RefreshExpire time.Duration // 刷新过期时间
}

type AppConfig struct {
	Server ServerConfig
	Log    LogConfig
	JWT    JWTConfig
	Redis  RedisConfig
	OTEL   OTELConfig
}

// RedisConfig 定义了 Redis 的配置.
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type OTELConfig struct {
	Enable      bool   `mapstructure:"enable"`
	ServiceName string `mapstructure:"service_name"`
	Version     string `mapstructure:"version"`
	Endpoint    string `mapstructure:"endpoint"`
	Insecure    bool   `mapstructure:"insecure"`
}
