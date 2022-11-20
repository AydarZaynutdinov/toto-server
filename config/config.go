package config

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"

	"gopkg.in/yaml.v3"
)

var (
	cfg  *Config
	once sync.Once
)

type (
	Config struct {
		App        App         `yaml:"app"`
		Logger     Logger      `yaml:"logger"`
		DB         DBConfig    `yaml:"db"`
		Redis      RedisConfig `yaml:"redis"`
		Migrations Migrations  `yaml:"migrations"`
	}

	App struct {
		Name        string `yaml:"name"`
		Port        string `yaml:"port" env:"PORT"`
		Host        string `yaml:"host" env:"HOST"`
		Description string `yaml:"description"`
		URI         string `yaml:"uri"`
		Environment string `yaml:"environment"`
		Version     string `yaml:"version"`
	}

	Logger struct {
		Level      string `yaml:"level"`
		TimeShow   bool   `yaml:"time_show"`
		TimeFormat string `yaml:"time_format"`
	}

	DBConfig struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     string `yaml:"port" env:"DB_PORT"`
		Username string `yaml:"username" env:"DB_USERNAME"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		Database string `yaml:"database" env:"DB_DATABASE"`
	}

	RedisConfig struct {
		Address string `yaml:"address" env:"REDIS_ADDRESS"`
	}

	Migrations struct {
		Dir string `yaml:"dir"`
	}
)

func New(configPath string) (*Config, error) {
	var err error
	once.Do(func() {
		cfg, err = parse(configPath)
	})

	return cfg, err
}

// Parses file by received parameter (filePath) to create app config
func parse(filePath string) (*Config, error) {
	filename, _ := filepath.Abs(filePath)
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	if err = cleanenv.ReadConfig(filename, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
