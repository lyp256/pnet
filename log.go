package pnet

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func SetLog(l *logrus.Logger) {
	logger = l
}
