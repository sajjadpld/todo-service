package config

type Service struct {
	Locale
	Env         string `mapstructure:"APP_ENV"`
	Name        string `mapstructure:"APP_NAME"`
	Debug       bool   `mapstructure:"APP_DEBUG"`
	LogLevel    string `mapstructure:"APP_LOG_LEVEL"`
	TimeZone    string `mapstructure:"APP_TIMEZONE"`
	StopTimeout int    `mapstructure:"APP_STOP_TIMEOUT"`
}
