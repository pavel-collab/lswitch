package translator

import (
    "errors"
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
// Возвращает ошибку, если конфиг найден, но некорректен.
func LoadConfig(path string) (Config, error) {
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

        // Валидируем конфигурацию из найденного файла
        if err := validateConfig(cfg); err != nil {
            return cfg, err
        }

        return cfg, nil
	}

    // Конфиг не найден — возвращаем значения по умолчанию без ошибки
    return cfg, nil
}

// getConfigCandidates возвращает список путей к конфигурационным файлам для проверки
func getConfigCandidates(userPath string) []string {
    candidates := []string{}

    // 1) Явно указанный путь через аргумент CLI имеет наивысший приоритет
    if userPath != "" {
        candidates = append(candidates, userPath)
    }

    // 2) Путь из переменной окружения LSWITCH_CONFIG (если задан)
    if envPath := os.Getenv("LSWITCH_CONFIG"); envPath != "" {
        candidates = append(candidates, envPath)
    }

    // 3) Локальный файл в текущей директории
    candidates = append(candidates, "./lswitch.yaml")

    // 4) XDG_CONFIG_HOME или ~/.config
    if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
        candidates = append(candidates, filepath.Join(xdg, "lswitch", "lswitch.yaml"))
    } else if home, err := os.UserHomeDir(); err == nil {
        candidates = append(candidates, filepath.Join(home, ".config", "lswitch", "lswitch.yaml"))
    }

    // 5) Файл в домашней директории пользователя (исторический вариант)
    if home, err := os.UserHomeDir(); err == nil {
        candidates = append(candidates, filepath.Join(home, ".lswitch.yaml"))
    }

    return candidates
}

// validateConfig проверяет корректность полей конфига
func validateConfig(cfg Config) error {
    switch strings.ToLower(strings.TrimSpace(cfg.Direction)) {
    case "en2ru", "ru2en", "":
        // пустое значение уже нормализовано ранее в en2ru
    default:
        return errors.New("invalid direction: must be 'en2ru' or 'ru2en'")
    }
    return nil
}
