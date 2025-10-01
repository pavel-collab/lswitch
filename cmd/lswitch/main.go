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

// version indicates the application version (injected at build time)
var version = "dev"

func main() {
    cfgPath := flag.String("config", "", "Path to YAML config file (optional)")
    showVersion := flag.Bool("version", false, "Show version and exit")
    showHelp := flag.Bool("help", false, "Show help and exit")
    useClipboard := flag.Bool("clipboard", false, "Read from and write to system clipboard")
    quiet := flag.Bool("quiet", false, "Do not print output to stdout in clipboard mode")
	flag.Parse()

    // Fast path: show help and exit
    if *showHelp {
        flag.Usage()
        os.Exit(0)
    }

	if *showVersion {
		fmt.Printf("lswitch version %s\n", version)
		os.Exit(0)
	}

    // Load configuration
    cfg, err := translator.LoadConfig(*cfgPath)
    if err != nil {
        log.Printf("invalid config: %v", err)
        os.Exit(2)
    }

    // Clipboard mode
    if *useClipboard {
        if err := handleClipboardMode(cfg, *quiet); err != nil {
            fmt.Fprintln(os.Stderr, "Clipboard mode error:", err)
            os.Exit(3)
        }
        return
    }

    // Read entire stdin
    input, err := readAllFromStdin()
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to read stdin:", err)
        os.Exit(3)
    }
    // If stdin is empty and clipboard mode is not selected — show usage and exit
    if input == "" {
        flag.Usage()
        return
    }

    // Create translator and transform text
	trans := translator.New(cfg)
	output := trans.Translate(string(input))

    // Write result to stdout
    _, err = io.WriteString(os.Stdout, output)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write stdout:", err)
        os.Exit(4)
	}
}

// readAllFromStdin reads all data from stdin. If stdin is not a pipe/file,
// it returns an empty string without error so the utility can no-op gracefully.
func readAllFromStdin() (string, error) {
    stat, err := os.Stdin.Stat()
    if err != nil {
        return "", err
    }
    // If no data is coming (stdin is a TTY), return an empty string
    if (stat.Mode() & os.ModeCharDevice) != 0 {
        return "", nil
    }
    b, err := io.ReadAll(os.Stdin)
    if err != nil {
        return "", err
    }
    return string(b), nil
}

// handleClipboardMode reads text from the system clipboard, translates it,
// writes it back, and optionally prints to stdout when quiet=false.
func handleClipboardMode(cfg translator.Config, quiet bool) error {
    text, err := clipboard.ReadAll()
    if err != nil {
        return fmt.Errorf("read clipboard: %w", err)
    }
    if text == "" {
        // Empty clipboard — do nothing
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
