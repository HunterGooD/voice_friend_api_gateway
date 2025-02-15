package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		CertFilePath string `yaml:"cert_file_path"`
	}
	Server struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Timeout struct {
			Server time.Duration `yaml:"server"`
			Write  time.Duration `yaml:"write"`
			Read   time.Duration `yaml:"read"`
			Idle   time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`
	Client struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"client"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBIdx    int    `yaml:"dbIdx"`
	} `yaml:"redis"`
}

func NewConfig(fileName string) (*Config, error) {
	var config Config

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

func (c *Config) BuildRedisConn() string {
	return fmt.Sprintf("redis://%s:%s@%s:%s/%d",
		c.Redis.User,
		c.Redis.Password,
		c.Redis.Host,
		c.Redis.Port,
		c.Redis.DBIdx,
	)
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
