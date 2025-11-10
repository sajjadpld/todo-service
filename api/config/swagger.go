package config

type Swagger struct {
	Host        string `mapstructure:"SWAGGER_HOST"`
	Schemes     string `mapstructure:"SWAGGER_SCHEMES"`
	Enable      bool   `mapstructure:"SWAGGER_ENABLE"`
	Title       string `mapstructure:"SWAGGER_INFO_TITLE"`
	Description string `mapstructure:"SWAGGER_INFO_DESCRIPTION"`
	Version     string `mapstructure:"SWAGGER_INFO_VERSION"`
	Username    string `mapstructure:"SWAGGER_USERNAME"`
	Password    string `mapstructure:"SWAGGER_PASSWORD"`
}
