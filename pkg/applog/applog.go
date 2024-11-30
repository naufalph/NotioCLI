package applog

import (
	"fmt"
	"log"
)

const (
	DebugLevel = "DEBUG"
	ErrorLevel = "ERROR"
	InfoLevel  = "INFO"
	WarnLevel  = "WARN"
	FatalLevel = "FATAL"
)

func cusomLog(logLevel string, err error, errorMsg string) {
	if err != nil {
		log.Println(generateErrMsgFormal(logLevel, err, errorMsg))
	} else {
		log.Println(generateLogMsgFormal(logLevel, errorMsg))
	}
}

func Debug(errorMsg string) {
	cusomLog(DebugLevel, nil, errorMsg)
}

func Error(err error, errorMsg string) {
	cusomLog(ErrorLevel, err, errorMsg)
}

func Info(err error, errorMsg string) {
	cusomLog(InfoLevel, err, errorMsg)
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
