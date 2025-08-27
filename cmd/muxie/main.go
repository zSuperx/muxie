// Package main provides main  
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/phanorcoll/muxie/internal/config"
	applog "github.com/phanorcoll/muxie/internal/log"
	"github.com/phanorcoll/muxie/internal/tui"
)

var (
	version = "dev"
	date    = "unknown"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug logging")
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("muxie version %s, built at %s\n", version, date)
		os.Exit(0)
	}

	if *debug {
		log.Println("debug mode enabled")
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		logDir := filepath.Join(home, ".config", "muxie")
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatal(err)
		}

		logFile := filepath.Join(logDir, "debug.log")
		f, err := tea.LogToFile(logFile, "debug")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}

	logger := applog.New(*debug)

	config, err := config.Load()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	if len(config.Sessions) > 0 {
		logger.Printf("configuration loaded successfully")
	} else {
		logger.Printf("no sessions found in config, starting with default")
	}

	m := tui.NewModel(config, logger, version)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Printf("error starting: %v", err)
		fmt.Printf("Upss, there's been an error: %v", err)
		os.Exit(1)
	}
}
