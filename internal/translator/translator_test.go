package translator

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected Config
	}{
		{
			name: "default config",
			path: "",
			expected: Config{
				Direction: "en2ru",
				CustomMap: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := LoadConfig(tt.path)
			if cfg.Direction != tt.expected.Direction {
				t.Errorf("Expected direction %s, got %s", tt.expected.Direction, cfg.Direction)
			}
		})
	}
}

func TestTranslator_Translate(t *testing.T) {
	cfg := Config{Direction: "en2ru"}
	translator := New(cfg)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple en2ru",
			input:    "hello",
			expected: "руддщ",
		},
		{
			name:     "upper case",
			input:    "Hello",
			expected: "Руддщ",
		},
		{
			name:     "mixed case and symbols",
			input:    "Hello, World!",
			expected: "Руддщ, Щщдлз!",
		},
		{
			name:     "numbers unchanged",
			input:    "test123",
			expected: "еуе123",
		},
		{
			name:     "russian unchanged in en2ru",
			input:    "привет",
			expected: "привет",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translator.Translate(tt.input)
			if result != tt.expected {
				t.Errorf("Translate(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTranslator_TranslateRu2En(t *testing.T) {
	cfg := Config{Direction: "ru2en"}
	translator := New(cfg)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple ru2en",
			input:    "руддщ",
			expected: "hello",
		},
		{
			name:     "upper case ru2en",
			input:    "Руддщ",
			expected: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translator.Translate(tt.input)
			if result != tt.expected {
				t.Errorf("Translate(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBuildMapping(t *testing.T) {
	tests := []struct {
		name     string
		cfg      Config
		checkKey rune
		expected rune
	}{
		{
			name:     "en2ru mapping",
			cfg:      Config{Direction: "en2ru"},
			checkKey: 'q',
			expected: 'й',
		},
		{
			name:     "ru2en mapping",
			cfg:      Config{Direction: "ru2en"},
			checkKey: 'й',
			expected: 'q',
		},
		{
			name: "custom mapping",
			cfg: Config{
				Direction: "en2ru",
				CustomMap: map[string]string{"a": "ф", "b": "и"},
			},
			checkKey: 'a',
			expected: 'ф',
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping := buildMapping(tt.cfg)
			result, exists := mapping[tt.checkKey]
			if !exists {
				t.Errorf("Key %c not found in mapping", tt.checkKey)
			}
			if result != tt.expected {
				t.Errorf("Mapping(%c) = %c, expected %c", tt.checkKey, result, tt.expected)
			}
		})
	}
}

// BenchmarkTranslate бенчмарк для измерения производительности
func BenchmarkTranslate(b *testing.B) {
	cfg := Config{Direction: "en2ru"}
	translator := New(cfg)
	testString := "Hello, World! This is a test string for benchmarking."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		translator.Translate(testString)
	}
}
