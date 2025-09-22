package translator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config представляет конфигурацию транслятора
type Config struct {
	Direction string            `yaml:"direction"` // "en2ru" или "ru2en"
	CustomMap map[string]string `yaml:"custom_map,omitempty"`
}

// LoadConfig загружает конфигурацию из файла или использует значения по умолчанию
func LoadConfig(path string) Config {
	cfg := Config{
		Direction: "en2ru",
		CustomMap: nil,
	}

	candidates := getConfigCandidates(path)

	for _, p := range candidates {
		if p == "" {
			continue
		}

		data, err := os.ReadFile(p)
		if err != nil {
			continue // Файл не найден - пробуем следующий
		}

		var loadedConfig Config
		if err := yaml.Unmarshal(data, &loadedConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse config %s: %v\n", p, err)
			continue
		}

		// Обновляем конфигурацию
		if loadedConfig.Direction != "" {
			cfg.Direction = strings.ToLower(loadedConfig.Direction)
		}
		if len(loadedConfig.CustomMap) > 0 {
			cfg.CustomMap = loadedConfig.CustomMap
		}

		return cfg
	}

	return cfg
}

// getConfigCandidates возвращает список путей к конфигурационным файлам для проверки
func getConfigCandidates(userPath string) []string {
	candidates := []string{}

	if userPath != "" {
		candidates = append(candidates, userPath)
	}

	// Локальный файл в текущей директории
	candidates = append(candidates, "./layout_switcher.yaml")

	// Файл в домашней директории пользователя
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".layout_switcher.yaml"))
	}

	return candidates
}
