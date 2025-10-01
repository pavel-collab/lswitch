package translator

import (
	"strings"
	"unicode"
)

// Translator provides text transformation between keyboard layouts
type Translator struct {
	mapping map[rune]rune
}

// New creates a new translator instance with the given configuration
func New(cfg Config) *Translator {
	return &Translator{
		mapping: buildMapping(cfg),
	}
}

// Translate converts text according to the translator settings
func (t *Translator) Translate(s string) string {
    var builder strings.Builder
    // Pre-allocate approximately based on input byte length
    builder.Grow(len(s))

    for _, r := range s {
        // Try direct mapping first
        if to, ok := t.mapping[r]; ok {
            builder.WriteRune(to)
            continue
        }

        // Handle uppercase runes efficiently
        if isUpper := unicode.IsUpper(r); isUpper {
            lowerRune := unicode.ToLower(r)
            if to, ok := t.mapping[lowerRune]; ok {
                if unicode.IsLetter(to) {
                    builder.WriteRune(unicode.ToUpper(to))
                } else {
                    builder.WriteRune(to)
                }
                continue
            }
        }

        // Rune not found in mapping — keep as is
        builder.WriteRune(r)
    }

	return builder.String()
}

// buildMapping creates a rune translation table based on the configuration
func buildMapping(cfg Config) map[rune]rune {
	// Базовая таблица en -> ru
	baseMap := map[rune]rune{
		'q': 'й', 'w': 'ц', 'e': 'у', 'r': 'к', 't': 'е', 'y': 'н', 'u': 'г',
		'i': 'ш', 'o': 'щ', 'p': 'з', '[': 'х', ']': 'ъ', 'a': 'ф', 's': 'ы',
		'd': 'в', 'f': 'а', 'g': 'п', 'h': 'р', 'j': 'о', 'k': 'л', 'l': 'д',
		';': 'ж', '\'': 'э', 'z': 'я', 'x': 'ч', 'c': 'с', 'v': 'м', 'b': 'и',
		'n': 'т', 'm': 'ь', ',': 'б', '.': 'ю', '/': '.', '`': 'ё',
		// Символы с Shift
		'{': 'Х', '}': 'Ъ', ':': 'Ж', '"': 'Э', '<': 'Б', '>': 'Ю', '?': ',',
		'~': 'Ё', '$': ';', '^': ':', '&': '?',
	}

	// Цифры и некоторые символы остаются без изменений
	for r := '0'; r <= '9'; r++ {
		baseMap[r] = r
	}
	identicalSymbols := "!@#%*()_+-=|"
	for _, r := range identicalSymbols {
		baseMap[r] = r
	}

    // Apply user-provided mappings
	if cfg.CustomMap != nil {
		for key, value := range cfg.CustomMap {
			if key == "" || value == "" {
				continue
			}
			keyRunes := []rune(key)
			valueRunes := []rune(value)
			if len(keyRunes) == 1 && len(valueRunes) == 1 {
				baseMap[keyRunes[0]] = valueRunes[0]
			}
		}
	}

    // Invert mapping when direction is ru2en
    if cfg.Direction == "ru2en" {
		invertedMap := make(map[rune]rune, len(baseMap))
		for key, value := range baseMap {
			invertedMap[value] = key
		}
		return invertedMap
	}

	return baseMap
}
