package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func SetupFileLogging() error {
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

	// Set log output to both file and console
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	// Set log format with timestamp and file location
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return nil
}
