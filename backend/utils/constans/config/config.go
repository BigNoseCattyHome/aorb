package config

import (
	"github.com/spf13/viper"
	"os"
)

var Conf *Config

type Config struct {
	Server   *Server             `toml:"Server"`
	MongoDB  *MongoDB            `toml:"MongoDB"`
	Redis    *Redis              `toml:"Redis"`
	JWT      *JWT                `toml:"JWT"`
	Etcd     *Etcd               `toml:"Etcd"`
	Services map[string]*Service `toml:"Services"`
}

type Server struct {
	Port    string `toml:"port"`
	Version string `toml:"version"`
}

type MongoDB struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

type Redis struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Db       int    `toml:"db"`
}

type JWT struct {
	JwtSecret string `toml:"jwtSecret"`
}

type Etcd struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type Service struct {
	Name        string `toml:"name"`
	LoadBalance bool   `toml:"loadBalance"`
	Host        string `toml:"host"`
	Port        string `toml:"port"`
}

func InitConfig() {
	work, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(work + "/backend/utils/constans/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Conf)
	if err != nil {
		panic(err)
	}
}
