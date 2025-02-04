package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"

	"github.com/caarlos0/env"
)

// Config – объект конфигурации приложения.
type Config struct {
	// 	DatabaseURL  – dsn для подключения к БД.
	DatabaseURL string `yaml:"database_url"`
	TS          int64  `yaml:"ts"`
	Recipes     bool   `yaml:"recipes"`
	RecipesTS   int64  `yaml:"recipes_ts"`
	Code        string `yaml:"code"`
}

// Initialize – функция инициализации конфига.
func Initialize(configPath string) (*Config, error) {

	configFile, err := os.Open(configPath)
	var c = &Config{}
	err = yaml.NewDecoder(configFile).Decode(c)
	if err != nil {
		return nil, fmt.Errorf("parse yaml error: %w", err)
	}
	err = env.Parse(c)
	return c, err
}

func (c Config) SaveConfigToYAML(cfg Config, filename string) {
	// Преобразование структуры в YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("Ошибка при преобразовании в YAML: %v", err)
	}

	// Запись данных в файл
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatalf("Ошибка при записи в файл: %v", err)
	}

	fmt.Printf("Конфигурация успешно сохранена в %s\n", filename)
}
