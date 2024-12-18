package manager

import (
	"github.com/xiusin/pine/contracts"
)

type logger struct {
	contracts.Logger
}

func (l logger) Info(msg string, keysAndValues ...any) {
	l.Logger.Warn(msg, keysAndValues)
}

func (l logger) Error(err error, msg string, keysAndValues ...any) {
	l.Logger.Error("%s: 错误: %s, 参数: %s", msg, err, keysAndValues)
}
