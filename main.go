package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/pchchv/logsrv/htpasswd"
	_ "github.com/pchchv/logsrv/httpupstream"
	"github.com/pchchv/logsrv/logging"
	"github.com/pchchv/logsrv/login"
)

const appName = "logsrv"

func main() {
	config := login.ReadConfig()
	if err := logging.Set(config.LogLevel, config.TextLogging); err != nil {
		exit(nil, err)
	}
	logging.AccessLogCookiesBlacklist = append(logging.AccessLogCookiesBlacklist, config.CookieName)
	configToLog := *config
	configToLog.JwtSecret = "..."
	logging.LifecycleStart(appName, configToLog)
	h, err := login.NewHandler(config)
	if err != nil {
		exit(nil, err)
	}
	handlerChain := logging.NewLogMiddleware(h)
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	port := config.Port
	if port != "" {
		port = fmt.Sprintf(":%s", port)
	}
	httpSrv := &http.Server{Addr: port, Handler: handlerChain}
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				logging.ServerClosed(appName)
			} else {
				exit(nil, err)
			}
		}
	}()
	logging.LifecycleStop(appName, <-stop, nil)
	ctx, ctxCancel := context.WithTimeout(context.Background(), config.GracePeriod)
	httpSrv.Shutdown(ctx)
	ctxCancel()
}

var exit = func(signal os.Signal, err error) {
	logging.LifecycleStop(appName, signal, err)
	if err == nil {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
