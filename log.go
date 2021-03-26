package pnet

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// SetLog 设置日志
func SetLog(l *logrus.Logger) {
	logger = l
}
