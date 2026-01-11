package provider

import "github.com/spf13/viper"

type ENV string

const (
	ENV_DEVELOPMENT ENV = "development"
	ENV_TEST        ENV = "test"
	ENV_PRODUCTION  ENV = "production"
)

type Config struct {
	v        *viper.Viper
	App      AppConfig  `mapstructure:"app" json:"app" yaml:"app"`
	Database Database   `mapstructure:"database" json:"database" yaml:"database"`
	Redis    Redis      `mapstructure:"redis" json:"redis" yaml:"redis"`
	HTTP     HTTPConfig `mapstructure:"http" json:"http" yaml:"http"`
	Log      LogConfig  `mapstructure:"log" json:"log" yaml:"log"`
}

type AppConfig struct {
	Environment ENV    `mapstructure:"environment" json:"environment" yaml:"environment"`
	Timezone    string `mapstructure:"timezone" json:"timezone" yaml:"timezone"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"` // json, text
	Output     string `mapstructure:"output"` // stdout, file
	File       string `mapstructure:"file"`
	MaxSize    int    `mapstructure:"max_size"`    // MB
	MaxBackups int    `mapstructure:"max_backups"` // count
	MaxAge     int    `mapstructure:"max_age"`     // days
	Compress   bool   `mapstructure:"compress"`    // true/false

}

type Database struct {
	Host            string `mapstructure:"host" json:"host" yaml:"host"`
	Port            int    `mapstructure:"port" json:"port" yaml:"port"`
	User            string `mapstructure:"user" json:"user" yaml:"user"`
	Password        string `mapstructure:"password" json:"password" yaml:"password"`
	Name            string `mapstructure:"name" json:"name" yaml:"name"`
	Debug           bool   `mapstructure:"debug" json:"debug" yaml:"debug"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
}

type Redis struct {
	Address    string `mapstructure:"address" json:"address" yaml:"address"`
	DB         int    `mapstructure:"db" json:"db" yaml:"db"`
	Username   string `mapstructure:"username" json:"username" yaml:"username"`
	Password   string `mapstructure:"password" json:"password" yaml:"password"`
	TLSEnabled bool   `mapstructure:"tls_enabled" json:"tls_enabled" yaml:"tls_enabled"`
}

type HTTPConfig struct {
	Address      string `mapstructure:"host" json:"host" yaml:"host"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout" json:"write_timeout" yaml:"write_timeout"`
}

func NewConfig(v *viper.Viper) (*Config, error) {
	var c Config

	c.v = v
	c.setDefault()

	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) setDefault() {
	c.v.SetDefault("app.environment", "development")
	c.v.SetDefault("app.timezone", "Asia/Shanghai")

	// Database
	c.v.SetDefault("database.host", "localhost")
	c.v.SetDefault("database.port", 5432)
	c.v.SetDefault("database.user", "postgres")
	c.v.SetDefault("database.password", "postgres")
	c.v.SetDefault("database.name", "dev")
	c.v.SetDefault("database.debug", false)
	c.v.SetDefault("database.max_idle_conns", 10)
	c.v.SetDefault("database.max_open_conns", 100)
	c.v.SetDefault("database.conn_max_lifetime", 3600)

	// Redis
	c.v.SetDefault("redis.address", "localhost:6379")
	c.v.SetDefault("redis.db", 0)
	c.v.SetDefault("redis.username", "")
	c.v.SetDefault("redis.password", "")
	c.v.SetDefault("redis.tls_enabled", false)

	// HTTP
	c.v.SetDefault("http.host", "0.0.0.0")
	c.v.SetDefault("http.port", 8080)
	c.v.SetDefault("http.read_timeout", 10)
	c.v.SetDefault("http.write_timeout", 10)

	// Log
	c.v.SetDefault("log.level", "info")
	c.v.SetDefault("log.format", "json")
	c.v.SetDefault("log.output", "stdout")
	c.v.SetDefault("log.file", "logs/app.log")
	c.v.SetDefault("log.max_size", 100)
	c.v.SetDefault("log.max_backups", 7)
	c.v.SetDefault("log.max_age", 30)
	c.v.SetDefault("log.compress", false)
}
