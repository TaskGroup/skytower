package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	ApiAuth `yaml:"api_auth"`
}

type ApiAuth struct {
	Login     string `yaml:"login" env-required:"true"`
	Password  string `yaml:"password" env-required:"true"`
	Host      string `yaml:"host_url" env-required:"true"`
	SecretKey string `yaml:"secret_key"`
}

// #Must приставка ставится тогда, когда функция вместо возврата ошибки будет паниковать
func MustLoad() *Config {
	const configPath = "config/local.yaml"

	pathToBack, err := os.Executable()
	if err != nil {
		log.Fatalf("Config file error path: %s", err)
	}
	index := strings.LastIndex(pathToBack, "/back/")
	if index == -1 {
		fmt.Println("Подстрока не найдена для конфигурационного файла")
	}
	index += 5
	backPath := pathToBack[:index]
	pathToConfig := filepath.Join(backPath, configPath)
	if _, err = os.Stat(pathToConfig); os.IsNotExist(err) {
		log.Fatalf("Config file does not exists: %s", pathToConfig)
	}

	var cfg Config
	if err = cleanenv.ReadConfig(pathToConfig, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
