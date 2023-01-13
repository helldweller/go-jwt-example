package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/sync/errgroup"

	"package/main/internal/config"
	"package/main/internal/controllers"
	docs "package/main/internal/docs"
	"package/main/internal/middleware"
)

var ctx, cancel = context.WithCancel(context.Background())
var group, groupCtx = errgroup.WithContext(ctx)
var conf *config.Config

func init() {
	conf = config.Cfg
}

func Run() {

	log.Info("Starting app")

	gin.SetMode(gin.TestMode)
	// DebugMode
	// ReleaseMode
	// TestMode

	r := gin.New()

	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/metrics"}}))

	// @title 			GO API with JWT auth example
	// @version			1.0
	// @description		A example service API in Go using Gin framework.
	// @contact.name	Dmitry
	// @contact.url		https://github.com/helldweller
	// @contact.email	helldweller@protonmail.com
	// @license.name	MIT
	// @license.url		https://www.mit.edu/~amini/LICENSE.md
	// @host			localhost:8080
	// @BasePath		/api/v1
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/signup", controllers.Signup)
			r.GET("/validate", middleware.RequireAuth, controllers.Validate)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	server := &http.Server{
		Addr:    conf.HTTPListenIPPort,
		Handler: r,
		// BaseContext: ctx,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	group.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		defer close(signalChannel)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
		select {
		case sig := <-signalChannel:
			log.Errorf("Received signal: %s", sig)
			if err := server.Shutdown(ctx); err != nil {
				log.Errorf("Received an error while shutting down the server: %s", err)
			}
			cancel()
		case <-groupCtx.Done():
			log.Error("Closing signal goroutine")
			if err := server.Shutdown(ctx); err != nil {
				log.Errorf("Received an error while shutting down the server: %s", err)
			}
			return groupCtx.Err()
		}
		return nil
	})

	group.Go(func() error {
		log.Infof("Starting web server on %s", conf.HTTPListenIPPort)
		err := server.ListenAndServe()
		return err
	})

	err := group.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Error("Context was canceled")
		} else {
			log.Errorf("Received error: %v\n", err)
		}
	} else {
		log.Error("Sucsessfull finished")
	}
}
