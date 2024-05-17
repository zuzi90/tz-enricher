package logger

import "github.com/sirupsen/logrus"

func NewLogger(lvl string) (*logrus.Logger, error) {
	log := logrus.New()
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		level = logrus.DebugLevel
		log.Info("set default logging level: debug", lvl)
	}

	log.SetLevel(level)
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return log, nil
}
