package guthi_log

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var glog *log.Logger

var outputType string

type Levels string

const (
	L_Trace Levels = "L_TRACE"
	L_Debug Levels = "L_DEBUG"
	L_Info  Levels = "L_INFO"
	L_Warn  Levels = "L_WARN"
	L_Error Levels = "L_ERROR"
	L_Fatal Levels = "L_FATAL"
	L_Panic Levels = "L_PANIC"
)

type Formats string

const (
	F_TEXT Formats = "F_TEXT"
	F_JSON Formats = "F_JSON"
)

type Outputs string

const (
	O_STDOUT Outputs = "O_STDOUT"
	O_STDERR Outputs = "O_STDERR"
	O_FILE   Outputs = "O_FILE"
)

type LogType string

const (
	APPLICATION LogType = "APPLICATION"
	SYSTEM      LogType = "SYSTEM"
	ACCESS      LogType = "ACCESS"
)

var applicationLogFile *os.File
var systemLogFile *os.File
var accessLogFile *os.File

// var logLevel log.Level
var err error

func InitLog(level Levels, format Formats, output Outputs) {
	glog = log.New()

	// Set the log level
	switch level {
	case "L_TRACE":
		glog.Level = log.TraceLevel
	case "L_DEBUG":
		glog.Level = log.DebugLevel
	case "L_INFO":
		glog.Level = log.InfoLevel
	case "L_WARN":
		glog.Level = log.WarnLevel
	case "L_ERROR":
		glog.Level = log.ErrorLevel
	case "L_FATAL":
		glog.Level = log.FatalLevel
	case "L_PANIC":
		glog.Level = log.PanicLevel
	default:
		glog.Level = log.InfoLevel
	}

	// Set the log format
	switch format {
	case "F_TEXT":
		glog.Formatter = &log.TextFormatter{}
	case "F_JSON":
		glog.Formatter = &log.JSONFormatter{}
	default:
		glog.Formatter = &log.TextFormatter{}
	}

	// Set the log output
	switch output {
	case "O_STDOUT":
		outputType = "O_STDOUT"
		glog.Out = os.Stdout
	case "O_STDERR":
		outputType = "O_STDERR"
		glog.Out = os.Stderr
	case "O_FILE":
		outputType = "O_FILE"
		applicationLogFile, err = os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		systemLogFile, err = os.OpenFile("system.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		accessLogFile, err = os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		// glog.Out = systemLogFile
		// glog.Out = accessLogFile
		// glog.Out = applicationLogFile

	default:
		glog.Out = os.Stdout
	}

}

func Info(message string, logType string, fields map[string]interface{}) {
	if outputType == "O_FILE" {
		switch logType {
		case "APPLICATION":
			glog.Out = applicationLogFile
		case "SYSTEM":
			glog.Out = systemLogFile
		case "ACCESS":
			glog.Out = accessLogFile
		default:
			glog.Out = applicationLogFile
		}

	}
	glog.WithFields(fields).Info(message)
}
