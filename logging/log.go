package logging

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func Logger(app ...string) (logger *log.Logger) {
	logger = log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetReportCaller(true)
	logger.SetFormatter(&log.TextFormatter{})
	//log.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	return
}
