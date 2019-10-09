package models

import (
	"fmt"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func InitLogger() {
	log = logrus.New()
	log.SetLevel(logrus.DebugLevel)

	log.SetReportCaller(true)
	log.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}
}
