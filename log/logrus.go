package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.TraceLevel)
	Logger.SetReportCaller(true)

	//Logger.SetFormatter(&logrus.TextFormatter{
	//	ForceQuote:      true,
	//	TimestampFormat: "2006-01-02 15:04:05",
	//	FullTimestamp:   true,
	//})
	Logger.SetFormatter(&MyFormatter{})
}

func InitLogger() {
	Logger.SetLevel(LogrusCfg.Level)
	Logger.Info("Succeed to init log!")
}
