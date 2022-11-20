package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"toto-server/config"
	"toto-server/internal/consts"
	"toto-server/internal/repository/sql"
	"toto-server/internal/service"

	"github.com/fasthttp/router"
	"github.com/go-redis/redis/v9"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func Run(
	ctx context.Context,
	cancelFunc context.CancelFunc,
	appConfig config.Config,
	zapLog *zap.Logger,
	sqlDb *sql.SQL,
	redisClient *redis.Client,
) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// ------ sql repository ------
	skuConfigSqlRepo := sql.New(
		sql.WithContext(ctx),
		sql.WithConnection(sqlDb))

	// ------ service ------
	skuConfigService := service.New(
		service.WithLogger(zapLog),
		service.WithSkuConfigRepo(skuConfigSqlRepo),
		service.WithRedisClient(redisClient))

	// ------ server
	server := prepareHttpServer(ctx, skuConfigService)

	// ------ start http server ------
	go func() {
		port, err := strconv.Atoi(appConfig.App.Port)
		if err != nil {
			zapLog.Sugar().Fatal(map[string]interface{}{
				consts.Module: consts.AppModule,
				consts.Action: "parse http server port",
				consts.Params: fmt.Sprintf("port: %s", appConfig.App.Port),
				consts.Error:  err.Error(),
			})
		}
		if err = server.ListenAndServe(fmt.Sprintf(":%d", port)); err != nil {
			zapLog.Sugar().Fatal(map[string]interface{}{
				consts.Module: consts.AppModule,
				consts.Action: "http server listen",
				consts.Params: fmt.Sprintf("address: %s:%s", appConfig.App.Host, appConfig.App.Port),
				consts.Error:  err.Error(),
			})
		}
	}()

	zapLog.Sugar().Info(map[string]interface{}{
		consts.Module:  consts.AppModule,
		consts.Action:  "service started",
		consts.Version: appConfig.App.Version,
		consts.Port:    appConfig.App.Port,
	})

	sig := <-sigChan
	zapLog.Sugar().Info(map[string]interface{}{
		consts.Module: consts.AppModule,
		consts.Action: fmt.Sprintf("start graceful shutdown, caught sig: %+v", sig),
	})

	zapLog.Sugar().Info(map[string]interface{}{
		consts.Module: consts.AppModule,
		consts.Action: "shutdown server is done",
	})

	// ------ stop service context ------
	cancelFunc()

	zapLog.Sugar().Info(map[string]interface{}{
		consts.Module: consts.AppModule,
		consts.Action: "shutdown processes of the service are done",
	})
	os.Exit(0)
}

func prepareHttpServer(ctx context.Context, skuConfigService service.IService) *fasthttp.Server {
	r := router.New()
	r.PanicHandler = func(ctx *fasthttp.RequestCtx, i interface{}) {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusInternalServerError))
	}
	r.NotFound = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusNotFound))
	}
	r.MethodNotAllowed = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed))
	}

	r.GET(consts.MetricRequestURI, fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))
	r.GET(consts.GetMainSkuURI, skuConfigService.Get(ctx))

	server := &fasthttp.Server{
		Handler: r.Handler,
	}
	return server
}
