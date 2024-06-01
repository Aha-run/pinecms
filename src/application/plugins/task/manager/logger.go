package manager

import (
	pineLogger "github.com/xiusin/logger"
)

type logger struct {
	pineLogger.AbstractLogger
}

func (l logger) Info(msg string, keysAndValues ...any) {
	l.AbstractLogger.Print(msg, keysAndValues)
}

func (l logger) Error(err error, msg string, keysAndValues ...any) {
	l.Errorf("%s: 错误: %s, 参数: %s", msg, err, keysAndValues)
}
