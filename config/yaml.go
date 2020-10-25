package config

import (
	"fmt"

	"github.com/toolkits/pkg/file"
)

type Config struct {
	Logger       loggerSection   `yaml:"logger"`
	Alidyvms     AlidyvmsSection `yaml:"alidyvms"`
	Consumer     consumerSection `yaml:"consumer"`
	Redis        redisSection    `yaml:"redis"`
	MaxDelayTime int             `yaml:"maxDelayTime"`
}

type loggerSection struct {
	Dir       string `yaml:"dir"`
	Level     string `yaml:"level"`
	KeepHours uint   `yaml:"keepHours"`
}

type redisSection struct {
	Addr    string         `yaml:"addr"`
	Pass    string         `yaml:"pass"`
	Idle    int            `yaml:"idle"`
	Timeout timeoutSection `yaml:"timeout"`
}

type timeoutSection struct {
	Conn  int `yaml:"conn"`
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
}

type consumerSection struct {
	Queue  string `yaml:"queue"`
	Worker int    `yaml:"worker"`
}

type AlidyvmsSection struct {
	AccessKey        string   `yaml:"access_key"`
	Secret           string   `yaml:"secret"`
	RegionId         string   `yaml:"region_id"`
	TtsCode          string   `yaml:"tts_code"`
	CalledShowNumber []string `yaml:"called_numbers"`
}

var yaml Config

func Get() Config {
	return yaml
}

func ParseConfig(yf string) error {
	err := file.ReadYaml(yf, &yaml)
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", yf, err)
	}
	return nil
}
