package config

import "github.com/spf13/viper"

type DBConfig struct {
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
}

func LoadDBConfig() (config DBConfig, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("db")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
