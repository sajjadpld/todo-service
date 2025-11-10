package config

import "time"

type Http struct {
	Host         string        `mapstructure:"HTTP_HOST"`
	Port         string        `mapstructure:"HTTP_PORT"`
	WriteTimeout time.Duration `mapstructure:"HTTP_WRITE_TIMEOUT"`
	ReadTimeout  time.Duration `mapstructure:"HTTP_READ_TIMEOUT"`
	Tls          bool          `mapstructure:"HTTP_TLS"`
}
