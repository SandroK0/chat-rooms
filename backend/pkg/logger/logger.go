package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func InitLogger() error {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return err
	}

	// Create or open the log file
	logFile, err := os.OpenFile(
		filepath.Join(logsDir, "app.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		return err
	}

	// Create multi-writers for both file and console output
	infoWriter := io.MultiWriter(os.Stdout, logFile)
	warningWriter := io.MultiWriter(os.Stdout, logFile)
	errorWriter := io.MultiWriter(os.Stderr, logFile)

	// Initialize loggers with different prefixes
	InfoLogger = log.New(infoWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(warningWriter, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

// Convenience functions
func Info(v ...interface{}) {
	InfoLogger.Println(v...)
}

func Infof(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

func Warning(v ...interface{}) {
	WarningLogger.Println(v...)
}

func Warningf(format string, v ...interface{}) {
	WarningLogger.Printf(format, v...)
}

func Error(v ...interface{}) {
	ErrorLogger.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}