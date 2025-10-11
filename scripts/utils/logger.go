package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logFile *os.File
var logger *log.Logger

func InitLogger() error {
	var err error
	logFile, err = os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logger = log.New(logFile, "", log.LstdFlags)
	LogInfo("=== Installation script started ===")
	return nil
}

func CloseLogger() {
	if logFile != nil {
		LogInfo("=== Installation script finished ===")
		logFile.Close()
	}
}

func LogInfo(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[INFO] %s: %s", timestamp, message)

	// Print to console
	fmt.Println(logEntry)

	// Write to log file
	if logger != nil {
		logger.Println(logEntry)
	}
}

func LogError(message string, err error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[ERROR] %s: %s - %v", timestamp, message, err)

	// Print to console
	fmt.Println(logEntry)

	// Write to log file
	if logger != nil {
		logger.Println(logEntry)
	}
}

func LogCommand(command string) {
	LogInfo(fmt.Sprintf("Executing command: %s", command))
}