package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logFilePath = fmt.Sprintf("./%v.log", time.Now().Format("2006-01-02"))
var logger *log.Logger = nil

func Printf(format string, v ...interface{}) {
	if logger == nil {
		if logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666); err != nil {
			fmt.Println(err)
		} else {
			logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
		}
	}
	logger.Printf(format, v...)
}
