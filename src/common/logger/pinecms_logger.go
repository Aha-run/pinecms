package logger

import (
	"bytes"
	"strings"
	"time"

	"github.com/xiusin/logger"
	"github.com/xiusin/pine"
	"github.com/xiusin/pinecms/src/application/models/tables"
	"github.com/xiusin/pinecms/src/common/helper"
	"xorm.io/xorm"
)

var firstLineChar = "pine.(*routerWrapper).result"
var _sp = "\n"

type pineCmsLoggerWriter struct {
	orm       *xorm.Engine
	logCh     chan []byte
	closed    bool
	errorFlag []byte
}

func NewPineCmsLogger(orm *xorm.Engine, len uint) *pineCmsLoggerWriter {
	l := &pineCmsLoggerWriter{orm: orm, logCh: make(chan []byte, len), errorFlag: []byte("[ERRO]")}
	go l.BeginConsume()
	return l
}

func (p *pineCmsLoggerWriter) BeginConsume() {
	for {
		log, isCloser := <-p.logCh
		if !isCloser {
			return
		}
		if _, err := p.orm.InsertOne(p.parseLog(log)); err != nil {
			pine.Logger().Warning("日志入库失败", err)
		}
	}
}

func (p *pineCmsLoggerWriter) Write(data []byte) (int, error) {
	defer func() {
		if err := recover(); err != nil {
			p.closed = true
		}
	}()
	if !p.closed && bytes.Contains(data, p.errorFlag) {
		p.logCh <- data
	}
	return 0, nil
}

func (p *pineCmsLoggerWriter) parseLog(log []byte) *tables.Log {
	lines := strings.Split(*helper.Bytes2String(log), _sp)
	var index = -1
	for i, v := range lines {
		if strings.Contains(v, firstLineChar) {
			index = i
			break
		}
	}
	if index-3 > 3 {
		lines = append(lines[0:1], lines[3:index-3]...)
	}
	return &tables.Log{Level: uint8(logger.ErrorLevel), Message: strings.Join(lines, _sp), Time: tables.LocalTime(time.Now())}
}

func (p *pineCmsLoggerWriter) Close() {
	close(p.logCh)
}
