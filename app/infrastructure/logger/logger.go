package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger(fields logrus.Fields) *logrus.Entry {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.TraceLevel)
	logger.ReportCaller = true
	logger.SetOutput(os.Stdout)

	entry := logger.WithFields(fields)

	// closeFn := func() {file.Close()}

	return entry
}
