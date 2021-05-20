package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"coinswallet/internal/services"
	"coinswallet/internal/transport"
	"github.com/go-kit/kit/log"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8081", "HTTP listen address")
		s        services.Service
	)

	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	s = services.NewService(logger)
	h := transport.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))

	errs := make(chan error)
	go func() {
		_ = logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()
	go handleSignals(errs)

	_ = logger.Log("Exit", <-errs)
}

func handleSignals(c chan error) {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	c <- fmt.Errorf("%s", <-sigCh)
}
