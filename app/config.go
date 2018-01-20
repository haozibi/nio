package app

import (
	"github.com/haozibi/nio/utils"

	yaml "gopkg.in/yaml.v2"
)

type NioConfig struct {
	Client struct {
		ServerIP   string `yaml:"server_ip"`
		ServerPort string `yaml:"server_port"`
	}
	Server struct {
		BindIP   string `yaml:"bind_ip"`
		BindPort string `yaml:"bind_port"`
	}
	App []struct {
		Name       string `yaml:"name"`
		LocalPort  string `yaml:"local_port"`
		BindIP     string `yaml:"bind_ip"`
		ListenPort string `yaml:"listen_port"`
		Passwd     string `yaml:"passwd"`
	}
	Log struct {
		LogLevel  string `yaml:"log_level"`
		LogPath   string `yaml:"log_path"`
		LogWay    string `yaml:"log_way"`
		LogSimple bool   `yaml:"log_simple"`
	}
}

var CONF = new(NioConfig)

func init() {
	body := utils.Readfile("./conf.yaml")
	err := yaml.Unmarshal(body, CONF)
	if err != nil {
		panic(err)
	}
}
