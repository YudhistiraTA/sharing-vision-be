package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YudhistiraTA/sharing-vision-be/internal/logger"
	"github.com/YudhistiraTA/sharing-vision-be/internal/repository/mysql"
	"github.com/YudhistiraTA/sharing-vision-be/internal/rest"
	"github.com/YudhistiraTA/sharing-vision-be/model"
	"github.com/YudhistiraTA/sharing-vision-be/service/article"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// logger init
	log := logger.DefaultLogger()

	// env init
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to read .env", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		}, zapcore.Field{
			Key:    "env",
			Type:   zapcore.StringType,
			String: ".env",
		})
	}

	// runtime ctx
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	//prepare database
	dbDriver := os.Getenv("DATABASE_DRIVER")
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatal("Failed to open connection to database", zapcore.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		})
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal("Failed to ping database ", zapcore.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		}, zapcore.Field{
			Key:    "connection",
			Type:   zapcore.StringType,
			String: connection,
		})
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal("Failed to close DB connection", zapcore.Field{
				Key:       logger.ErrorField,
				Type:      zapcore.ErrorType,
				Interface: err,
			})
		}
	}()

	// prepare server
	sc := model.ServerConfig{
		Addr: os.Getenv("PORT"),
		Db:   dbConn,
		Log:  log,
	}

	// prepare services
	ar := mysql.NewArticleRepository(sc.Db)
	as := article.NewArticleService(ar)

	// prepare server
	srv := rest.NewServer(sc, as)

	// run group
	rungroup, ctx := errgroup.WithContext(ctx)
	rungroup.Go(func() error {
		if er := srv.ListenAndServe(); er != nil && !errors.Is(er, http.ErrServerClosed) {
			return fmt.Errorf("listen and server %w", er)
		}
		return nil
	})
	rungroup.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if er := srv.Shutdown(shutdownCtx); er != nil {
			return fmt.Errorf("shutdown http server %w", er)
		}

		return nil
	})
	if err := rungroup.Wait(); err != nil {
		log.Error("run group exited because of error", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		})
		return
	}
	log.Info("server exited properly")
}
