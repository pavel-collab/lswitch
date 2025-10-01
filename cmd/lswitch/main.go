package main

import (
	"flag"
	"fmt"
	"io"
    "log"
	"os"

	"lswitch/internal/translator"
    clipboard "github.com/atotto/clipboard"
)

// version указывает версию приложения (задаётся при сборке)
var version = "dev"

func main() {
    cfgPath := flag.String("config", "", "Path to YAML config file (optional)")
    showVersion := flag.Bool("version", false, "Show version and exit")
    showHelp := flag.Bool("help", false, "Show help and exit")
    useClipboard := flag.Bool("clipboard", false, "Read from and write to system clipboard")
    quiet := flag.Bool("quiet", false, "Do not print output to stdout in clipboard mode")
	flag.Parse()

    // Быстрый путь: справка
    if *showHelp {
        flag.Usage()
        os.Exit(0)
    }

	if *showVersion {
		fmt.Printf("lswitch version %s\n", version)
		os.Exit(0)
	}

    // Загружаем конфигурацию
    cfg, err := translator.LoadConfig(*cfgPath)
    if err != nil {
        log.Printf("invalid config: %v", err)
        os.Exit(2)
    }

    // Режим работы с буфером обмена
    if *useClipboard {
        if err := handleClipboardMode(cfg, *quiet); err != nil {
            fmt.Fprintln(os.Stderr, "Clipboard mode error:", err)
            os.Exit(3)
        }
        return
    }

    // Читаем stdin полностью
    input, err := readAllFromStdin()
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to read stdin:", err)
        os.Exit(3)
    }
    // Если stdin пуст и не выбран режим clipboard — показать usage и выйти
    if input == "" {
        flag.Usage()
        return
    }

	// Создаём транслятор и преобразуем текст
	trans := translator.New(cfg)
	output := trans.Translate(string(input))

	// Записываем результат в stdout
    _, err = io.WriteString(os.Stdout, output)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write stdout:", err)
        os.Exit(4)
	}
}

// readAllFromStdin читает все данные из stdin. Если stdin не подключен к пайпу/файлу,
// функция возвращает пустую строку без ошибки, чтобы утилита могла просто ничего не делать.
func readAllFromStdin() (string, error) {
    stat, err := os.Stdin.Stat()
    if err != nil {
        return "", err
    }
    // Если данные не поступают (stdin — терминал), просто вернуть пустую строку
    if (stat.Mode() & os.ModeCharDevice) != 0 {
        return "", nil
    }
    b, err := io.ReadAll(os.Stdin)
    if err != nil {
        return "", err
    }
    return string(b), nil
}

// handleClipboardMode читает текст из системного буфера обмена, переводит его и
// записывает обратно. Опционально печатает результат в stdout, если quiet=false.
func handleClipboardMode(cfg translator.Config, quiet bool) error {
    text, err := clipboard.ReadAll()
    if err != nil {
        return fmt.Errorf("read clipboard: %w", err)
    }
    if text == "" {
        // Пустой буфер обмена — ничего не делаем
        return nil
    }
    trans := translator.New(cfg)
    out := trans.Translate(text)
    if err := clipboard.WriteAll(out); err != nil {
        return fmt.Errorf("write clipboard: %w", err)
    }
    if !quiet {
        _, _ = io.WriteString(os.Stdout, out)
    }
    return nil
}
