package registry

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	src "microservice"
	"log"
)

type registry struct {
	*viper.Viper
}

func New() IRegistry {
	return &registry{
		viper.New(),
	}
}

func (v *registry) Init() {
	v.AddConfigPath(src.Root())
	v.SetConfigName(EnvFormat)
	v.SetConfigType(EnvMime)
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("[registry] init failure: ", err)
	}
}

func (v *registry) Parse(item interface{}) {
	err := v.Unmarshal(&item)
	if err != nil {
		log.Fatal(err)
	}
}
