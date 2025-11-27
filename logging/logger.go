package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Logger struct {
	logLevel     string
	logToFile    bool
	logFilePath  string
	outputFormat string
	logFile      *os.File
}

func NewLogger(
	level string,
	logToFile bool,
	filePath string,
	outputFormat string,
) *Logger {
	logger := &Logger{
		logLevel:     level,
		logToFile:    logToFile,
		outputFormat: strings.ToLower(outputFormat),
	}

	if logToFile {
		ext := ".log"
		if strings.ToLower(outputFormat) == "json" {
			ext = ".json"
		}

		base := strings.TrimSuffix(filePath, ext)
		timestamp := time.Now().Format("2006-01-02_15-04")
		dynamicFilePath := fmt.Sprintf("%s_%s%s", base, timestamp, ext)

		// Make sure the directory exists
		dir := filepath.Dir(dynamicFilePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if mkErr := os.MkdirAll(dir, 0755); mkErr != nil {
				fmt.Printf("Warning: Could not create directory %s: %v\n", dir, mkErr)
			}
		}

		file, err := os.OpenFile(dynamicFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Warning: Could not open log file %s: %v\n", dynamicFilePath, err)
		} else {
			logger.logFilePath = dynamicFilePath
			logger.logFile = file
		}
	}

	return logger
}

func (l *Logger) LogInfo(message string)  { l.log("INFO", message) }
func (l *Logger) LogDebug(message string) { l.log("DEBUG", message) }
func (l *Logger) LogWarn(message string)  { l.log("WARN", message) }
func (l *Logger) LogError(message string) { l.log("ERROR", message) }

func (l *Logger) LogTerminal(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	logMessage := fmt.Sprintf("[%s] %s", timestamp, message)
	if l.logToFile && l.logFile != nil {
		fmt.Fprintln(l.logFile, logMessage)
	}
	fmt.Println(message)
}

func (l *Logger) shouldLog(requiredLevel string) bool {
	levels := map[string]int{
		"debug": 4,
		"info":  3,
		"warn":  2,
		"error": 1,
	}
	requiredLevelIndex := levels[requiredLevel]
	currentLevelIndex := levels[l.logLevel]
	return currentLevelIndex >= requiredLevelIndex
}

func (l *Logger) log(level, message string) {
	if !l.shouldLog(strings.ToLower(level)) {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	var logMessage string

	if l.outputFormat == "json" {
		entry := map[string]interface{}{
			"timestamp": timestamp,
			"level":     level,
			"message":   message,
		}
		jsonBytes, _ := json.Marshal(entry)
		logMessage = string(jsonBytes)
	} else {
		logMessage = fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)
		if strings.Contains(message, "=== STAGE") {
			logMessage = fmt.Sprintf("\n%s\n[%s] [%s] %s\n%s\n",
				strings.Repeat("=", 50),
				timestamp,
				level,
				message,
				strings.Repeat("=", 50),
			)
		}
	}

	if l.logToFile && l.logFile != nil {
		fmt.Fprintln(l.logFile, logMessage)
	} else {
		fmt.Println(logMessage)
	}
}
