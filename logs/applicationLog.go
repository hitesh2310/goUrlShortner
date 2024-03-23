package logs

import (
	"fmt"
	"io"
	"main/pkg/constants"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var ApplicationLog *logrus.Logger

func SetUpApplicationLogs() {
	logger := logrus.New()

	// Set the log level (debug, info, warn, error, fatal, panic)
	logger.SetLevel(logrus.DebugLevel)
	// fmt.Sprintf(format, a)
	// You can set different output formats. JSON is a common choice.
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Optionally, you can add more outputs. Here, we are adding stdout and a log file.
	// logger.SetOutput(os.Stdout)

	// currentTimestamp := time.Now().Format("2006-01-02") // YYYY-MM-DD format

	// Dynamically generate the log file name with the current timestamp
	logFileName := constants.ApplicationConfig.Application.LogPath + "goBulkCampaignService.log"
	// file, err := os.OpenFile("./logs/logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	logFile := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    100,
		MaxAge:     28,
		MaxBackups: 30000,
		LocalTime:  true,
		Compress:   true,
	}

	logger.SetOutput(io.MultiWriter(logFile))
	ApplicationLog = logger

	// log.Println("Failed to log to file, using default stderr")
	defer logFile.Close()

}

func InfoLog(format string, a ...any) {
	stringMessage := fmt.Sprintf(format, a...)
	ApplicationLog.WithFields(logrus.Fields{}).Info(stringMessage)
}

func ErrorLog(format string, a ...any) {
	stringMessage := fmt.Sprintf(format, a...)
	ApplicationLog.WithFields(logrus.Fields{}).Error(stringMessage)
}
