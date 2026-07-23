//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package boot

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/algotiqa/core"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/spf13/viper"
)

//=============================================================================
//===
//=== Public functions
//===
//=============================================================================

func ReadConfig(component string, config any) {
	viper.SetConfigName(component)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/algotiqa/")
	viper.AddConfigPath("$HOME/.algotiqa/" + component)
	viper.AddConfigPath("config")

	//--- Use env vars to replace config

	viper.AutomaticEnv()
	viper.SetEnvPrefix("ALGO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//--- Watch and reload the configuration

	viper.WatchConfig()

	//--- Read config

	err := viper.ReadInConfig()
	core.ExitIfError(err)

	err = viper.Unmarshal(config)
	core.ExitIfError(err)
}

//=============================================================================

func InitLogger(component, version string, app *core.Application) *slog.Logger {
	//--- Create log file

	logFile := "log/" + component + ".log"

	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	core.ExitIfError(err)

	var wrt io.Writer = f

	if !app.Production {
		wrt = io.MultiWriter(os.Stdout, f)
	}

	//--- create logger

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	if !app.Debug {
		opts = nil
	}

	logger := slog.New(slog.NewJSONHandler(wrt, opts)).With(
		slog.String("component", component),
		slog.Int("pid", os.Getpid()),
	)

	slog.SetDefault(logger)
	logStartupBanner(component, version)

	return logger
}

//=============================================================================

func InitEngine(logger *slog.Logger, app *core.Application) *gin.Engine {
	engine := gin.New()
	engine.Use(sloggin.New(logger))
	engine.Use(gin.Recovery())

	if app.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	return engine
}

//=============================================================================

func RunHttpServer(router *gin.Engine, app *core.Application) {
	slog.Info("Starting HTTPS server...")
	rootCAs, err := x509.SystemCertPool()
	core.ExitIfError(err)

	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	caCert, err := os.ReadFile("config/ca.crt")
	core.ExitIfError(err)

	if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
		core.ExitWithMessage("Failed to append CA cert to local certificate pool")
	}

	tlsConfig := &tls.Config{
		ClientCAs:  rootCAs,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr:      app.BindAddress,
		TLSConfig: tlsConfig,
		Handler:   router,
	}

	slog.Info("Running")
	err = server.ListenAndServeTLS("config/server.crt", "config/server.key")
	core.ExitIfError(err)
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func logStartupBanner(component, version string) {
	mem := core.GetMemoryInfo()

	slog.Info("=== Starting service ========================",
		"version",    version,
		"cpus",       runtime.NumCPU(),
		"usedMemory", fmt.Sprintf("%d MB", mem.UsedMB),
		"totalMemory",fmt.Sprintf("%d MB", mem.TotalMB),
	)
}

//=============================================================================
