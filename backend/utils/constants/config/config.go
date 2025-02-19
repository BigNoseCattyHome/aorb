package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Conf *Config

// 使用方法：
// config.Conf.MongoDB.Host

type Config struct {
	Consul   *Consul             `toml:"Consul"`
	Server   *Server             `toml:"Server"`
	MongoDB  *MongoDB            `toml:"MongoDB"`
	Redis    *Redis              `toml:"Redis"`
	JWT      *JWT                `toml:"JWT"`
	Services map[string]*Service `toml:"Services"`
	Pod      *Pod                `toml:"Pod"`
	Log      *Log                `toml:"Log"`
	Tracing  *Tracing            `toml:"Tracing"`
	RabbitMQ *RabbitMQ           `toml:"RabbitMQ"`
	Gorse    *Gorse              `toml:"Gorse"`
	Other    *Other              `toml:"Other"`
	// Etcd     *Etcd               `toml:"Etcd"`
	// PyroScope *PyroScope          `toml:"PyroScope"`
}

// 使用viper读取配置文件
func init() {
	work, _ := os.Getwd()
	fmt.Println("Current working directory:", work)
	fmt.Println("Config path:", work+"/../../utils/constants/config")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(work + "/../../utils/constants/config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Fatal error config file: %s \n", err)
		panic(err)
	}
	err = viper.Unmarshal(&Conf) // 将配置文件解析到Conf结构体中
	if err != nil {
		panic(err)
	}
}

type Other struct {
	AnonymityUser string `toml:"anonymityUser"`
}

type PyroScope struct {
	State string `toml:"state"`
	Addr  string `toml:"addr"`
}

type Gorse struct {
	GorseAddr   string `toml:"gorseAddr"`
	GorseApikey string `toml:"gorseApikey"`
}

type Tracing struct {
	EndPoint string  `toml:"endPoint"`
	State    string  `toml:"state"`
	Sampler  float64 `toml:"sampler"`
}

type RabbitMQ struct {
	Username    string `toml:"username"`
	Password    string `toml:"password"`
	Host        string `toml:"host"`
	Port        string `toml:"port"`
	VhostPrefix string `toml:"vhostPrefix"`
}

type Log struct {
	LoggerLevel string `toml:"loggerLevel"`
	LogPath     string `toml:"logPath"`
}

type Consul struct {
	Addr          string `toml:"addr"`
	AnonymityName string `toml:"anonymityName"`
}

type Pod struct {
	PodIp string `toml:"podIp"`
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
	Prefix   string `toml:"prefix"`
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
