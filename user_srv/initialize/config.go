package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	"mbSrvs/user_srv/global"
)

func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

func InitConfig() {
	data := GetEnvInfo("ENV")
	var configFileName string
	configFileNamePrefix := "config"
	configFileName = fmt.Sprintf("user_srv/%s-pro.yaml", configFileNamePrefix)
	if data == "local" {
		configFileName = fmt.Sprintf("user_srv/%s-debug.yaml", configFileNamePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
}
