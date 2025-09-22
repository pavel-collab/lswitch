package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// Config структура для YAML-конфига
type Config struct {
	Direction string            `yaml:"direction"` // "en2ru" или "ru2en"
	CustomMap map[string]string `yaml:"custom_map,omitempty"`
}

func main() {
	cfgPath := flag.String("config", "", "path to YAML config file (optional)")
	flag.Parse()

	cfg := loadConfig(*cfgPath)

	// читаем stdin полностью
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read stdin:", err)
		os.Exit(2)
	}

	m := buildMapping(cfg)

	out := translate(string(in), m)

	_, err = io.WriteString(os.Stdout, out)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to write stdout:", err)
		os.Exit(3)
	}

}

func loadConfig(path string) Config { // значения по умолчанию
	cfg := Config{Direction: "en2ru", CustomMap: nil}

	candidates := []string{}
	if path != "" {
		candidates = append(candidates, path)
	}
	// local file
	candidates = append(candidates, "./layout_switcher.yaml")
	// home
	home, err := os.UserHomeDir()
	if err == nil {
		candidates = append(candidates, filepath.Join(home, ".layout_switcher.yaml"))
	}

	for _, p := range candidates {
		if p == "" {
			continue
		}
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var got Config
		if err := yaml.Unmarshal(b, &got); err != nil {
			fmt.Fprintln(os.Stderr, "warning: failed to parse config", p, ":", err)
			continue
		}
		// если направление задано — заменяем
		if got.Direction != "" {
			cfg.Direction = strings.ToLower(got.Direction)
		}
		if len(got.CustomMap) > 0 {
			cfg.CustomMap = got.CustomMap
		}
		return cfg
	}

	return cfg

}

func buildMapping(cfg Config) map[rune]rune { // базовая таблица en -> ru
	base := map[rune]rune{'q': 'й', 'w': 'ц', 'e': 'у', 'r': 'к', 't': 'е', 'y': 'н', 'u': 'г', 'i': 'ш', 'o': 'щ', 'p': 'з', '[': 'х', ']': 'ъ', 'a': 'ф', 's': 'ы', 'd': 'в', 'f': 'а', 'g': 'п', 'h': 'р', 'j': 'о', 'k': 'л', 'l': 'д', ';': 'ж', '\'': 'э', 'z': 'я', 'x': 'ч', 'c': 'с', 'v': 'м', 'b': 'и', 'n': 'т', 'm': 'ь', ',': 'б', '.': 'ю', '/': '.', '`': 'ё',
		// shifted punctuation -> upper-case Russian letters or symbols
		'{': 'Х', '}': 'Ъ', ':': 'Ж', '"': 'Э', '<': 'Б', '>': 'Ю', '?': ',', '~': 'Ё'}

	// дополнительно сопоставим цифры и некоторые символы как идентичные (можно не менять)
	for r := '0'; r <= '9'; r++ {
		base[r] = r
	}

	// применим custom_map сверху (переводим ключи/значения в односимвольные rune)
	if cfg.CustomMap != nil {
		for k, v := range cfg.CustomMap {
			if k == "" || v == "" {
				continue
			}
			kr := []rune(k)
			vr := []rune(v)
			if len(kr) == 1 && len(vr) == 1 {
				base[kr[0]] = vr[0]
			}
		}
	}

	// если направление ru2en — инвертируем
	if cfg.Direction == "ru2en" {
		inv := make(map[rune]rune, len(base))
		for k, v := range base {
			inv[v] = k
		}
		return inv
	}

	return base

}

func translate(s string, m map[rune]rune) string {
	var b strings.Builder

	for _, r := range s {
		// пробуем прямое совпадение
		if to, ok := m[r]; ok {
			b.WriteRune(to)
			continue
		}
		// если буква верхнего регистра — попробуем понижать до нижнего, переводить, затем возвращать в верхний
		if unicode.IsUpper(r) {
			low := unicode.ToLower(r)
			if to, ok := m[low]; ok {
				// если результирующий символ буква — делаем верхний регистр
				if unicode.IsLetter(to) {
					b.WriteRune(unicode.ToUpper(to))
				} else {
					b.WriteRune(to)
				}

				continue
			}
		}
		// ничего не найдено — оставляем как есть
		b.WriteRune(r)
	}
	return b.String()
}
