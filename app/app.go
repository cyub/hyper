package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_config "github.com/cyub/hyper/config"
	"github.com/cyub/hyper/logger"
	config "github.com/cyub/hyper/pkg/config"
	"github.com/cyub/hyper/router"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// App struct
type App struct {
	*Options
	Gin               *gin.Engine
	Config            *config.Config
	Logger            *logrus.Logger
	Middlewares       []gin.HandlerFunc
	Components        []string
	ComponentMounters []ComponentMount
	RouterRegisters   []func(r *gin.Engine)
	Once              sync.Once
}

// ComponentMount use for mount component interface
type ComponentMount func(*App) error

// NewApp return an pointer to App
func NewApp(opts ...Option) *App {
	app := &App{
		Options: &Options{
			Name:          "hyper",
			Addr:          ":8000",
			RunMode:       "dev",
			CfgCenterAddr: "localhost:8500",
			CfgCenterPath: "",
		},
	}
	for _, opt := range opts {
		opt(app.Options)
	}
	return app
}

// RegisterRouter use for register router
func (app *App) RegisterRouter(f ...func(r *gin.Engine)) *App {
	app.RouterRegisters = append(app.RouterRegisters, f...)
	return app
}

// BootstrapWith use for register component
func (app *App) BootstrapWith(mounters ...ComponentMount) *App {
	app.ComponentMounters = append(app.ComponentMounters, mounters...)
	return app
}

// Bootstrap use for boot core components
func (app *App) Bootstrap() *App {
	app.Once.Do(func() {
		app.printBanner()
		app.bootConfig()
		app.bootBuildIns()
		app.bootLogger()
		app.bootMiddlewares()
		app.bootRoutes()
		app.bootComponents()
	})

	return app
}

// Use for add global middleware
func (app *App) Use(middlewares ...gin.HandlerFunc) *App {
	app.Middlewares = append(app.Middlewares, middlewares...)
	return app
}

// SetGin for inject gin Engin into App
func (app *App) SetGin(gin *gin.Engine) *App {
	app.Gin = gin
	return app
}

func (app *App) bootConfig() (err error) {
	if len(app.CfgCenterPath) == 0 {
		app.CfgCenterPath = fmt.Sprintf("%s/%s/config", app.Name, app.RunMode)
	}
	err = _config.Init(app.CfgCenterAddr, app.CfgCenterPath)
	if err != nil {
		panic(err)
	}
	app.Config = _config.Instance()
	if len(app.Addr) == 0 {
		app.Addr = app.Config.GetString("app.addr", ":8000")
	}
	switch app.RunMode {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "test", "staging":
		gin.SetMode(gin.TestMode)
	case "prod", "production":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
	return
}

func (app *App) bootBuildIns() {
	if app.Gin == nil {
		app.SetGin(gin.New())
	}

	if app.Config.GetBool("app.profile", true) {
		app.RegisterRouter(router.ProfileProvider())
	}
	if app.Config.GetBool("app.gin_logger", false) {
		app.Use(gin.Logger())
	}
}

func (app *App) bootLogger() (err error) {
	err = logger.Init(app.Config.GetString("log.writer", "stdout"),
		app.Config.GetString("log.level", "DEBUG"),
		app.Config.GetString("log.file", ""),
		app.Config.GetString("log.format", "text"))
	if err != nil {
		panic(err)
	}
	app.Logger = logger.Instance()
	return
}

func (app *App) bootRoutes() {
	for _, r := range app.RouterRegisters {
		r(app.Gin)
	}
}

func (app *App) bootComponents() {
	for _, component := range app.ComponentMounters {
		if err := component(app); err != nil {
			panic(err)
		}
	}
}

func (app *App) bootMiddlewares() {
	app.Gin.Use(app.Middlewares...)
}

// Run use for run application
func (app *App) Run() (err error) {
	app.Bootstrap()
	app.Logger.Infof("app name[%s] runmode[%s] addr[%s] run", app.Name, app.RunMode, app.Addr)
	srv := &http.Server{
		Addr:    app.Addr,
		Handler: app.Gin,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatalf("http listen fatal error %s\n", err)
		}
	}()

	app.WaitGracefulExit(srv)
	return
}

func (app *App) printBanner() {
	const banner = `
    __  __
   / / / / __  __  ____   ___    _____
  / /_/ / / / / / / __ \ / _ \  / ___/
 / __  / / /_/ / / /_/ //  __/ / /
/_/ /_/  \__, / / .___/ \___/ /_/
        /____/ /_/
Welcome to use Hyper, application is starting ...
`
	fmt.Printf("\x1b[32m%s\x1b[0m\n", banner)
}

// WaitGracefulExit graceful exit
func (app *App) WaitGracefulExit(srv *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			logger.Debug("server will exit")
			srv.Close()
			logger.Debug("server exited")
			return
		case syscall.SIGHUP:
		default:
		}
	}
}
