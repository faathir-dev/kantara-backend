package config

import (
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

const (
	LogDirectory = "./logs/query_log"
)

func SetupLogger() logger.Interface {
	// 1. Jika di Vercel / Production, gunakan stdout bawaan Go (tanpa buat file di disk)
	if os.Getenv("APP_ENV") == "production" || os.Getenv("VERCEL") == "1" || os.Getenv("IS_LOGGER") == "false" {
		return logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)
	}

	// 2. Lingkungan Lokal (tetap simpan log ke file)
	err := os.MkdirAll(LogDirectory, os.ModePerm)
	if err != nil {
		log.Printf("Warning: failed to create log directory: %v. Fallback to Stdout.", err)
		return logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)
	}

	currentMonth := strings.ToLower(time.Now().Format("January"))
	logFileName := currentMonth + "_query.log"

	logFile, err := os.OpenFile(LogDirectory+"/"+logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Warning: failed to open log file: %v. Fallback to Stdout.", err)
		return logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)
	}

	return logger.New(
		log.New(logFile, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      false,
		},
	)
}