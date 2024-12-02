package main

import (
	"os"
	"time"

	"backup-restorer/internal/model"
	"backup-restorer/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
)

func main() {
	// Настройка логирования
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Создание папки logs, если она не существует
	if err := os.MkdirAll("logs", 0755); err != nil {
		logrus.Fatalf("Failed to create logs directory: %v", err)
	}

	// Открытие файла для записи логов
	logFile, err := os.OpenFile("logs/backup-restorer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Настройка logrus для записи в файл
	logrus.SetOutput(logFile)

	p := tea.NewProgram(model.InitialModel())
	if err := p.Start(); err != nil {
		logrus.Errorf("Error: %v", err)
		time.Sleep(3 * time.Second) // Задержка в 3 секунды
		utils.ClearTerminal()
		os.Exit(1)
	}
	time.Sleep(3 * time.Second) // Задержка в 3 секунды
	utils.ClearTerminal()
}
