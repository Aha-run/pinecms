package taglibs

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/xiusin/pine"
	"github.com/xiusin/pinecms/src/application/controllers"
	"github.com/xiusin/pinecms/src/application/models/tables"
	"github.com/xiusin/pinecms/src/common/helper"
	"xorm.io/xorm"
)

var defaultArrReturnVal = reflect.ValueOf([]any{})

var defaultSignalVal = reflect.ValueOf(nil)

func checkArgType(args *jet.Arguments) bool {
	l := args.NumOfArguments()
	for i := 0; i < l; i++ {
		t := args.Get(i)
		if !isNumber(t) && t.Type().String() != "string" && t.Type().String() != "bool" {
			pine.Logger().Errorf("参数类型不支持: idx: %d -> type: %s -> val: %s", i, t, args.Get(i))
			return false
		}
	}
	return true
}

func isNumber(val reflect.Value) bool {
	return strings.Contains(val.String(), "float") || strings.Contains(val.String(), "int")
}

func getNumber(val reflect.Value) int64 {
	t := val.Type().String()
	if strings.Contains(t, "float") {
		return int64(val.Float())
	} else if strings.Contains(t, "int") {
		return val.Int()
	} else if t == "string" {
		v, _ := strconv.Atoi(val.String())
		return int64(v)
	}
	return 0
}

func getOrmSess(table ...string) *xorm.Session {
	if len(table) == 0 {
		table = []string{"articles"}
	}
	return helper.GetORM().Table(controllers.GetTableName(table[0]))
}

func getCategoryOrm() *xorm.Session {
	return helper.GetORM().Table(&tables.Category{}).Where("ismenu = 1")
}

func getCategoryTable() string {
	return controllers.GetTableName("category")
}

func getTagHash(args jet.Arguments) string {
	var arr []any
	for i := 0; i < args.NumOfArguments(); i++ {
		arr = append(arr, args.Get(i).Interface())
	}
	byts, _ := json.Marshal(&arr)
	md := md5.New()
	md.Write(byts)

	return fmt.Sprintf("%x", md.Sum(nil))
}
