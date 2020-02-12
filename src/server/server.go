package config

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/xiusin/iriscms/src/config"
	"github.com/xiusin/iriscms/src/router"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/boltdb"
	"github.com/kataras/iris/v12/view"
	"github.com/xiusin/iriscms/src/application/controllers"
	"github.com/xiusin/iriscms/src/application/controllers/backend"
	"github.com/xiusin/iriscms/src/application/controllers/middleware"
	"github.com/xiusin/iriscms/src/common/cache"
	"github.com/xiusin/iriscms/src/common/helper"
	"github.com/xiusin/iriscms/src/common/logger"
)

var (
	app    *iris.Application
	mvcApp *mvc.Application

	htmlEngine      *view.HTMLEngine
	XOrmEngine      *xorm.Engine
	sess            *sessions.Sessions
	sessionInitSync sync.Once
	sessCache       *boltdb.Database
	iCache          cache.ICache
	conf            = config.AppConfig()
)

func initDatabase() {
	dc := config.DBConfig()
	m, o := dc.Mysql, dc.Orm
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", m.DbUser, m.DbPassword, m.DbServer, m.DbPort, m.DbName, m.DbChatSet)
	_orm, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	_orm.SetLogger(logger.NewIrisCmsXormLogger(helper.NewOrmLogFile(conf.LogPath)))
	err = _orm.Ping() //检测是否联通数据库
	if err != nil {
		panic(err.Error())
	}
	_orm.ShowSQL(o.ShowSql)
	_orm.ShowExecTime(o.ShowExecTime)
	_orm.SetMaxOpenConns(int(o.MaxOpenConns))
	_orm.SetMaxIdleConns(int(o.MaxIdleConns))
	XOrmEngine = _orm
	//syncTable()
}

func initApp() {
	//实例化服务器
	app = iris.New()
	mvcApp = mvc.New(app).Configure(getMvcConfig())
	//配置前端缓存10秒
	//app.Use(iris.Cache304(10 * time.Second))
	//配置PPROF
	if conf.Pprof.Open {
		app.Logger().Debug("pprof enabled")
		app.Get(conf.Pprof.Route, pprof.New())
	}
	//附加视图
	djangoEngine := view.Django(conf.View.Path, ".tpl").Reload(conf.View.Reload)
	registerDjangoFunc(djangoEngine)
	app.RegisterView(djangoEngine)
	htmlEngine = view.HTML(conf.View.Path, ".html").Reload(conf.View.Reload)
	registerHtmlFunc(htmlEngine)
	app.RegisterView(htmlEngine)
}

func Server() {
	initDatabase()
	initApp()
	registerStatic()
	catchError()
	registerErrorRoutes()
	registerBackendRoutes()
	router.InitRouter(mvcApp)
	runServe()
}

func registerStatic() {
	app.Favicon(conf.Favicon, "favicon.ico")
	for _, static := range conf.Statics {
		app.HandleDir(static.Route, filepath.FromSlash(static.Path))
	}

}

//利用中间件执行控制器前置操作
func registerBackendRoutes() {
	mvcApp.Party(
		conf.BackendRouteParty,
		middleware.ViewRequestPath(app, conf.LogPath),
		middleware.CheckAdminLoginAndAccess(sess, iCache, XOrmEngine),
		middleware.SetGlobalConfigData(XOrmEngine, iCache),
		iris.Gzip,
	).Handle(new(backend.AdminController)).
		Handle(new(backend.LoginController)).
		Handle(new(backend.IndexController)).
		Handle(new(backend.CategoryController)).
		Handle(new(backend.ContentController)).
		Handle(new(backend.SettingController)).
		Handle(new(backend.SystemController)).
		Handle(new(backend.MemberController)).
		Handle(new(backend.DocumentController)).
		Handle(new(backend.LinkController))
	mvcApp.Party("/public", middleware.SetGlobalConfigData(XOrmEngine, iCache), injectConfig()).Handle(new(backend.PublicController))
}

//防止相互调用先用这种不优美的方式实现
func injectConfig() func(ctx context.Context) {
	return func(ctx context.Context) {
		ctx.Values().Set("app.config", iris.Map{"uploadEngine": conf.Upload.Engine})
		ctx.Next()
	}
}

func registerErrorRoutes() {
	err := new(backend.ErrorController)
	app.OnErrorCode(iris.StatusInternalServerError, err.ServerError)
	app.OnErrorCode(iris.StatusNotFound, err.StatusNotFound)
}

func runServe() {
	//golog.AddOutput(os.Stdout)
	if conf.Pprof.Open {
		go func() {
			pport := strconv.Itoa(int(conf.Pprof.Port))
			err := http.ListenAndServe(":"+pport, nil)
			if err != nil {
				app.Logger().Error("启动pprof失败", err)
			}
		}()
	}
	port := strconv.Itoa(int(conf.Port))
	_ = app.Run(iris.Addr(":"+port),
		iris.WithCharset(conf.Charset),
		iris.WithOptimizations,
		iris.WithPostMaxMemory(conf.Upload.MaxBodySize<<20),
	)
}

// 获取mvc配置, 分离相关session
func getMvcConfig() func(app *mvc.Application) {
	sessionInitSync.Do(func() {
		var err error
		hashKey, blockKey := []byte(conf.HashKey), []byte(conf.BlockKey)
		sec, ssc := securecookie.New(hashKey, blockKey), conf.Session
		sess = sessions.New(sessions.Config{Cookie: ssc.Name, Encode: sec.Encode, Decode: sec.Decode, Expires: ssc.Expires * time.Second})
		sessCache, err = boltdb.New(conf.CacheDb, os.FileMode(0750))
		if err != nil {
			app.Logger().Error("创建session缓存失败", err)
			panic(err)
		}
		iCache = cache.New(sessCache.Service, string(controllers.WebSiteCacheBucket))
		sess.UseDatabase(sessCache)
		iris.RegisterOnInterrupt(func() {
			if err := sessCache.Close(); err != nil {
				app.Logger().Error("关闭cache失败", err)
			}
		})
	})
	return func(m *mvc.Application) {
		// 注册依赖服务, 内部可以用反射类型方式获取
		m.Register(sess.Start, XOrmEngine, iCache)
	}
}

func getRequestLogs(ctx context.Context) string {
	var status, ip, method, path string
	status = strconv.Itoa(ctx.GetStatusCode())
	path = ctx.Path()
	method = ctx.Method()
	ip = ctx.RemoteAddr()
	return fmt.Sprintf("%v %s %s %s", status, path, method, ip)
}

func catchError() {
	app.Use(func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				if ctx.IsStopped() {
					return
				}
				var stacktrace string
				for i := 1; ; i++ {
					_, f, l, got := runtime.Caller(i)
					if !got {
						break
					}
					stacktrace += fmt.Sprintf("%s:%d\n", f, l)
				}
				logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.HandlerName())
				logMessage += fmt.Sprintf("At Request: %s\n", getRequestLogs(ctx))
				logMessage += fmt.Sprintf("Trace: %s\n", err)
				logMessage += fmt.Sprintf("\n%s", stacktrace)
				app.Logger().Error(logMessage)
				ctx.StatusCode(500)
				ctx.StopExecution()
			}
		}()
		ctx.Next()
	})
}
