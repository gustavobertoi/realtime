package log

import (
	"log"
	"os"

	"github.com/nolleh/caption_json_formatter"
	"github.com/sirupsen/logrus"
)

var _log = NewLogger()

func NewLogger() *logrus.Logger {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logrus.SetOutput(os.Stdout)
	logger := logrus.New()
	logger.Level = logrus.TraceLevel
	if os.Getenv("APP_DEBUG") == "1" {
		formatter := caption_json_formatter.Json()
		formatter.Colorize = true
		formatter.PrettyPrint = true
		logger.SetFormatter(formatter)
	}
	return logger
}

func GetStaticInstance() *logrus.Logger {
	return _log
}

func CreateWithContext(context string, fields logrus.Fields) *logrus.Entry {
	logger := NewLogger()
	fields["context"] = context
	return logger.WithFields(fields)
}
