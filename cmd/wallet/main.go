// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

// starts http server for the REST API of wallet

package main

import (
	"coinswallet/pkg/subprocmgr"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"coinswallet/internal/services"
	"coinswallet/internal/transport"
	"github.com/go-kit/kit/log"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8081", "HTTP listen address")

		s services.Service // services that implement business logic
	)

	// global program context
	ctx, cancel := context.WithCancel(context.Background())
	// goroutines manager
	goMgr := subprocmgr.New()

	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	s = services.NewService(logger)
	h := transport.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))

	// channel of "exit" signal
	errs := make(chan error)

	// run server
	goMgr.Add("httpServer")
	go func() {
		defer goMgr.Remove("httpServer")
		runHttpServer(ctx, &h, httpAddr, logger, errs)
	}()

	// check program termination
	go handleSignals(errs, func() {
		cancel()
	})

	// waiting for value in errs channel
	// it becomes when program terminated or http server shutdown
	_ = logger.Log("Exit", <-errs)

	// waiting for completing all goroutines
	<-goMgr.Done()
}

// runHttpServer - run http server and shutdown one correctly
func runHttpServer(ctx context.Context, h *http.Handler, httpAddr *string, logger log.Logger, errs chan error) {
	httpServer := &http.Server{
		Handler:      *h,
		Addr:         *httpAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	_ = logger.Log("transport", "HTTP", "addr", *httpAddr)
	stopHttpErrChan := make(chan error)
	go func() { stopHttpErrChan <- httpServer.ListenAndServe() }()

	select {
	case e := <-stopHttpErrChan:
		// http server error
		_ = logger.Log("httpServer", "terminate", "error", e)
		errs <- e
		return

	case <-ctx.Done():
		// exit program(global context "Done")
		shutdownCtx, c := context.WithTimeout(context.Background(), 5*time.Second) // set timeout 5sec for http server shutdown
		defer c()
		_ = logger.Log("http server", "shutdown", "result",
			httpServer.Shutdown(shutdownCtx),
		)
	}
}

// handleSignals - handle system interrupt signals and prepare program to finish
// the value in channel "c" is set and the function onExit is called
func handleSignals(c chan error, onExit func()) {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-sigCh
	onExit()
	c <- fmt.Errorf("%s", sig)
}
