package tplfun

import (
	"reflect"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/CloudyKit/jet"
	"github.com/xiusin/pine"
)

func format(str string) string {
	str = strings.Replace(str, "Y", "2006", 1)
	str = strings.Replace(str, "m", "01", 1)
	str = strings.Replace(str, "d", "02", 1)
	str = strings.Replace(str, "H", "13", 1)
	str = strings.Replace(str, "i", "04", 1)
	str = strings.Replace(str, "s", "05", 1)
	return str
}

func FormatTime(args jet.Arguments) reflect.Value {
	arg := strings.ReplaceAll("T", " ", args.Get(0).String())
	arg = strings.ReplaceAll("Z", "", arg)
	t, err := time.Parse(time.DateTime, arg)
	if err != nil {
		return reflect.ValueOf("")
	}
	format := time.DateTime
	if args.NumOfArguments() > 1 {
		format = args.Get(1).String()
	}
	return reflect.ValueOf(t.Format(format))
}

func MyDate(args jet.Arguments) reflect.Value {
	t, err := time.Parse("2006-01-02T15:04:05Z", args.Get(1).String())
	if err != nil {
		pine.Logger().Error("解析时间错误", err)
		return reflect.ValueOf("")
	}
	format := format(args.Get(0).String())
	return reflect.ValueOf(t.Format(format))
}

func GetDateTimeMK(args jet.Arguments) reflect.Value {
	return FormatTime(args)
}

func CnSubstr(args jet.Arguments) reflect.Value {
	me := args.Get(0).String()
	length := int(args.Get(1).Float())

	if utf8.RuneCountInString(me) > length {
		titleRune := []rune(me)
		me = string(titleRune[:length])
	}
	return reflect.ValueOf(me)
}
