package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	WatchInterval int    `json:"watch_interval"`
	DefaultTail   int    `json:"default_tail"`
	OutputFormat  string `json:"output_format"`
	WebhookURL    string `json:"webhook_url"`
	WebhookType   string `json:"webhook_type"` // "slack" o "discord"
}

var Defaults = Config{
	WatchInterval: 3,
	DefaultTail:   50,
	OutputFormat:  "table",
	WebhookURL:    "",
	WebhookType:   "",
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio home: %w", err)
	}
	return filepath.Join(home, ".wachiman", "config.json"), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return &Defaults, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Defaults, nil
		}
		return nil, fmt.Errorf("error leyendo config: %w", err)
	}

	cfg := Defaults
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config.json inválido: %w", err)
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("error creando directorio de config: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error guardando config: %w", err)
	}

	return nil
}

func Set(key, value string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	switch key {
	case "watch_interval":
		var v int
		if _, err := fmt.Sscanf(value, "%d", &v); err != nil {
			return fmt.Errorf("watch_interval debe ser un número entero")
		}
		if v < 1 {
			return fmt.Errorf("watch_interval mínimo es 1 segundo")
		}
		cfg.WatchInterval = v

	case "default_tail":
		var v int
		if _, err := fmt.Sscanf(value, "%d", &v); err != nil {
			return fmt.Errorf("default_tail debe ser un número entero")
		}
		if v < 1 {
			return fmt.Errorf("default_tail mínimo es 1")
		}
		cfg.DefaultTail = v

	case "output_format":
		if value != "table" && value != "json" {
			return fmt.Errorf("output_format debe ser 'table' o 'json'")
		}
		cfg.OutputFormat = value

	case "webhook_url":
		cfg.WebhookURL = value

	case "webhook_type":
		if value != "slack" && value != "discord" && value != "" {
			return fmt.Errorf("webhook_type debe ser 'slack' o 'discord'")
		}
		cfg.WebhookType = value

	default:
		return fmt.Errorf("campo desconocido: %q — opciones: watch_interval, default_tail, output_format, webhook_url, webhook_type", key)
	}

	return Save(cfg)
}