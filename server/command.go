package server

import (
	"context"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
)

var listenFlag = cli.StringFlag{
	Name:   "listen",
	Usage:  "Configure the server to listen to this interface.",
	EnvVar: "LISTEN",
	Value:  "0.0.0.0:5000",
}

var environmentFlag = cli.StringFlag{
	Name:   "environment",
	Usage:  "Set the environment the application is running in.",
	EnvVar: "ENVIRONMENT",
	Value:  "development",
}

var apiUserFlag = cli.StringFlag{
	Name:   "api-user",
	Usage:  "Set user of the api.",
	EnvVar: "API_USER",
	Value:  "aocp",
}

var apiPasswordFlag = cli.StringFlag{
	Name:   "api-password",
	Usage:  "Set password for the api user.",
	EnvVar: "API_PASSWORD",
	Value:  "aocp",
}

var storageFlag = cli.StringFlag{
	Name:   "storage",
	Usage:  "Set the storage engine.",
	EnvVar: "STORAGE",
	Value:  "local",
}

var storageArgsFlag = cli.StringFlag{
	Name:   "storage-args",
	Usage:  "Set the arguments used to connect to the storage engine.",
	EnvVar: "STORAGE_ARGS",
	Value:  "",
}

var Command = cli.Command{
	Name:  "server",
	Usage: "Run the server.",
	Flags: []cli.Flag{
		listenFlag,
		environmentFlag,
		apiUserFlag,
		apiPasswordFlag,
		storageFlag,
		storageArgsFlag,
	},
	Action: serverCommandAction,
}

func serverCommandAction(cliCtx *cli.Context) error {
	logger, err := getLogger(cliCtx)
	if err != nil {
		return err
	}

	logger.Info("Starting",
		zap.String("GOOS", runtime.GOOS),
		zap.String("env", cliCtx.String("environment")))

	store, err := getStorage(cliCtx)
	if err != nil {
		return err
	}

	h := handlers{
		logger:  logger,
		storage: store,
	}

	r := gin.New()

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	r.GET("/health", h.health)
	r.GET("/callback", h.callback)

	apiUser := cliCtx.String("api-user")
	apiPassword := cliCtx.String("api-password")

	enableAPI := true

	if len(apiUser) == 0 || len(apiPassword) == 0 {
		logger.Warn("API resources are disabled: user/password not set")
		enableAPI = false
	}

	if enableAPI {
		apiRouter := r.Group("/api")
		{
			accounts := make(gin.Accounts)
			accounts[apiUser] = apiPassword
			apiRouter.Use(gin.BasicAuth(accounts))

			apiRouter.POST("/locations", h.apiRecordLocation)
		}
	}

	var g run.Group

	srv := &http.Server{
		Addr:    getListenAddress(cliCtx),
		Handler: r,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g.Add(func() error {
		logger.Info("starting http service", zap.String("addr", srv.Addr))
		return srv.ListenAndServe()
	}, func(error) {
		logger.Info("stopping http service")
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("error shutting http service down", zap.Error(err))
		}
	})

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	g.Add(func() error {
		logger.Info("starting signal listener")
		<-quit
		return nil
	}, func(error) {
		logger.Info("stopping signal listener")
		close(quit)
	})

	if err := g.Run(); err != nil {
		logger.Error("error caught", zap.Error(err))
	}
	return nil
}

var ProductionEnvironment = "production"

func getLogger(c *cli.Context) (*zap.Logger, error) {
	if c.String("environment") == ProductionEnvironment {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

func getStorage(c *cli.Context) (storage, error) {
	switch c.String("storage") {
	case "local":
		return newLocalStorage()
	case "mysql":
		return newMySQLStorage(c)
	case "postgres":
		return newPGStorage(c)
	case "redis":
		return newRedisStorage(c)
	default:
		return nil, errStorageEngineNotFound
	}
}

func getListenAddress(cliCtx *cli.Context) string {
	if listen := cliCtx.String("listen"); len(listen) > 0 {
		return listen
	}
	return ":8080"
}
