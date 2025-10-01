package translator

import (
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "gopkg.in/yaml.v3"
)

// Config represents translator configuration
type Config struct {
    Direction string            `yaml:"direction"` // "en2ru" or "ru2en"
	CustomMap map[string]string `yaml:"custom_map,omitempty"`
}

// LoadConfig loads configuration from file or falls back to defaults.
// Returns an error when a found config is invalid.
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
            continue // File not found — try next candidate
		}

		var loadedConfig Config
		if err := yaml.Unmarshal(data, &loadedConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse config %s: %v\n", p, err)
            continue
		}

        // Update configuration from loaded values
		if loadedConfig.Direction != "" {
			cfg.Direction = strings.ToLower(loadedConfig.Direction)
		}
		if len(loadedConfig.CustomMap) > 0 {
			cfg.CustomMap = loadedConfig.CustomMap
		}

        // Validate configuration from the discovered file
        if err := validateConfig(cfg); err != nil {
            return cfg, err
        }

        return cfg, nil
	}

    // Config not found — return defaults without error
    return cfg, nil
}

// getConfigCandidates returns a prioritized list of config paths to check
func getConfigCandidates(userPath string) []string {
    candidates := []string{}

    // 1) Explicit path from CLI flag has the highest priority
    if userPath != "" {
        candidates = append(candidates, userPath)
    }

    // 2) Path from env var LSWITCH_CONFIG (if set)
    if envPath := os.Getenv("LSWITCH_CONFIG"); envPath != "" {
        candidates = append(candidates, envPath)
    }

    // 3) Local file in current directory
    candidates = append(candidates, "./lswitch.yaml")

    // 4) XDG_CONFIG_HOME or ~/.config
    if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
        candidates = append(candidates, filepath.Join(xdg, "lswitch", "lswitch.yaml"))
    } else if home, err := os.UserHomeDir(); err == nil {
        candidates = append(candidates, filepath.Join(home, ".config", "lswitch", "lswitch.yaml"))
    }

    // 5) Historical location in user's home directory
    if home, err := os.UserHomeDir(); err == nil {
        candidates = append(candidates, filepath.Join(home, ".lswitch.yaml"))
    }

    return candidates
}

// validateConfig verifies config fields are valid
func validateConfig(cfg Config) error {
    switch strings.ToLower(strings.TrimSpace(cfg.Direction)) {
    case "en2ru", "ru2en", "":
        // empty value has been normalized to en2ru earlier
    default:
        return errors.New("invalid direction: must be 'en2ru' or 'ru2en'")
    }
    return nil
}
