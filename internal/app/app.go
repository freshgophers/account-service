package app

import (
	"account-service/pkg/sms"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"account-service/internal/config"
	"account-service/internal/handler"
	"account-service/internal/repository"
	"account-service/internal/service/account"
	"account-service/internal/service/otp"
	"account-service/pkg/log"
	"account-service/pkg/server"
)

const (
	schema      = "accounts"
	version     = "1.0.0"
	description = "account-service"
)

// Run initializes whole application.
func Run() {
	logger := log.New(version, description)

	cfg, err := config.New()
	if err != nil {
		logger.Error("ERR_INIT_CONFIG", zap.Error(err))
		return
	}

	// Dependencies
	smsClient := sms.New(
		sms.Credentials{
			Endpoint: cfg.SMS.Endpoint,
			Username: cfg.SMS.Username,
			Password: cfg.SMS.Password,
		})

	// Initializations
	repositories, err := repository.New(
		repository.WithPostgresStore(schema, cfg.POSTGRES.DSN))
	if err != nil {
		logger.Error("ERR_INIT_REPOSITORY", zap.Error(err))
		return
	}
	defer repositories.Close()

	accountService, err := account.New(
		account.WithUsersRepository(repositories.User))
	if err != nil {
		logger.Error("ERR_INIT_ACCOUNT_SERVICE", zap.Error(err))
		return
	}

	otpService, err := otp.New(
		otp.WithSMSClient(smsClient),
		otp.WithSecretRepository(repositories.Secret),
		otp.WithAccountService(accountService),
	)
	if err != nil {
		logger.Error("ERR_INIT_OTP_SERVICE", zap.Error(err))
		return
	}

	handlers, err := handler.New(
		handler.Dependencies{
			Configs:        cfg,
			OTPService:     otpService,
			AccountService: accountService,
		},
		handler.WithHTTPHandler())
	if err != nil {
		logger.Error("ERR_INIT_HANDLER", zap.Error(err))
		return
	}

	servers, err := server.New(
		server.WithHTTPServer(handlers.HTTP, cfg.HTTP.Port))
	if err != nil {
		logger.Error("ERR_INIT_SERVER", zap.Error(err))
		return
	}

	// Run our server in a goroutine so that it doesn't block.
	if err = servers.Run(logger); err != nil {
		logger.Error("ERR_RUN_SERVER", zap.Error(err))
		return
	}

	// Graceful Shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the httpServer gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	quit := make(chan os.Signal, 1) // create channel to signify a signal being sent

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel
	<-quit                                             // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")

	// create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err = servers.Stop(ctx); err != nil {
		panic(err) // failure/timeout shutting down the httpServer gracefully
	}

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here

	fmt.Println("Server was successful shutdown.")
}
