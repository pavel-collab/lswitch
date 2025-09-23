package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"lswitch/internal/translator"
)

// version указывает версию приложения (задаётся при сборке)
var version = "dev"

func main() {
	cfgPath := flag.String("config", "", "Path to YAML config file (optional)")
	showVersion := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("lswitch version %s\n", version)
		os.Exit(0)
	}

	// Загружаем конфигурацию
	cfg := translator.LoadConfig(*cfgPath)

	// Читаем stdin полностью
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read stdin:", err)
		os.Exit(2)
	}

	// Создаём транслятор и преобразуем текст
	trans := translator.New(cfg)
	output := trans.Translate(string(input))

	// Записываем результат в stdout
	_, err = io.WriteString(os.Stdout, output)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write stdout:", err)
		os.Exit(3)
	}
}
