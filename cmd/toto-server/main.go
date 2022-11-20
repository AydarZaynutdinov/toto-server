package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"toto-server/config"
	"toto-server/internal/app"
	"toto-server/internal/consts"
	"toto-server/internal/repository/sql"

	"github.com/go-redis/redis/v9"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"go.uber.org/zap"
)

var (
	ctx        context.Context
	cancelFunc context.CancelFunc

	configFilePath string
	appConfig      *config.Config
	zapLogger      *zap.Logger
	sqlDb          *sql.SQL
)

func init() {
	flagConfig := flag.String("config", "config.yaml", "config file path")
	flag.Parse()

	configFilePath = *flagConfig
}

func main() {
	var (
		err error
	)

	ctx, cancelFunc = context.WithCancel(context.Background())

	appConfig, err = config.New(configFilePath)
	if err != nil {
		log.Fatalf("failed to get config instance, configFile: %s\nerr: %v", configFilePath, err)
	}

	getZapLogger()

	sqlDb, err = sql.ProvideSQLConnection(&appConfig.DB, ctx)
	if err != nil {
		log.Fatalf("failed to get connection with sql db\nerr: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{Addr: appConfig.Redis.Address})

	runMigrate()

	app.Run(ctx, cancelFunc, *appConfig, zapLogger, sqlDb, redisClient)
}

func getZapLogger() {
	var (
		err error
	)

	if appConfig.App.Environment == consts.ProdEnvironment {
		zapLogger, err = zap.NewProduction()
		if err != nil {
			log.Fatalln(err)
		}

		zapLogger.Info("production zap logger started")
	} else {
		zapLogger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatalln(err)
		}

		zapLogger.Info("dev zap logger started")
	}

	defer func(zapLogger *zap.Logger) {
		_ = zapLogger.Sync() //nolint
	}(zapLogger)
}

func runMigrate() {
	source := fmt.Sprintf("file://%s", appConfig.Migrations.Dir)
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		appConfig.DB.Username, appConfig.DB.Password, appConfig.DB.Host, appConfig.DB.Port, appConfig.DB.Database)
	migration, err := migrate.New(source, dbUrl)
	if err != nil {
		log.Fatalf("failed to prepare migration; err: %v", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to execute migrations\nerr: %v", err)
	}
	zapLogger.Info("migration successfully finished")
}
