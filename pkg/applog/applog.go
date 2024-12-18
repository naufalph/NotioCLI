package applog

import (
	"fmt"
	"log"
	"os"
)

const (
	DebugLevel = "DEBUG"
	ErrorLevel = "ERROR"
	InfoLevel  = "INFO"
	WarnLevel  = "WARN"
	FatalLevel = "FATAL"
)

func InitLog() {
	logFile, err := os.OpenFile("applog.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func customLog(logLevel string, err error, errorMsg string) {
	if err != nil {
		log.Println(generateErrMsgFormal(logLevel, err, errorMsg))
	} else {
		log.Println(generateLogMsgFormal(logLevel, errorMsg))
	}
}

func Debug(errorMsg string) {
	customLog(DebugLevel, nil, errorMsg)
}

func Error(err error, errorMsg string) {
	customLog(ErrorLevel, err, errorMsg)
	log.Fatalln()
}

func Info(err error, errorMsg string) {
	customLog(InfoLevel, err, errorMsg)
}

func generateErrMsgFormal(logLevel string, err error, errorMsg string) string {
	if errorMsg != "" {
		return fmt.Sprintf("[%s] %s \n %s",
			logLevel,
			errorMsg,
			err.Error())
	} else {
		return fmt.Sprintf("[%s] Unable to find errorMsg ! \n %s",
			logLevel,
			err.Error())
	}
}

func generateLogMsgFormal(logLevel string, errorMsg string) string {
	if errorMsg != "" {
		return fmt.Sprintf("[%s] %s",
			logLevel,
			errorMsg)
	} else {
		return fmt.Sprintf("[%s] Unable to find errorMsg !", logLevel)
	}
}
