package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ankur12345678/uptime-monitor/pkg/logger"
)

const (
	RUNNING = iota
	TERMINATING
)

// Atomic struct encapsulating server state
type ServerState struct {
	running int32
}

// Atomically store that the server is shutting down
func (state *ServerState) Shutdown() {
	atomic.StoreInt32(&state.running, TERMINATING)
}

// Atomically check that is running
func (state *ServerState) IsRunning() bool {
	return atomic.LoadInt32(&state.running) == RUNNING
}

// type Graceful will run an http.Server, with graceful shutdown
type Graceful struct {
	HTTPServer      *http.Server
	ShutdownTimeout time.Duration
	State           *ServerState
}

func optimiseListenAddress(addr string) string {
	if strings.HasPrefix(addr, ":") {
		return "localhost" + addr
	}
	return addr
}

// ListenAndServe the enclosed http.Server or grpc.Server, but shutdown gracefully
func (g *Graceful) ListenAndServe(startupMessages ...string) {
	for _, s := range startupMessages {
		logger.Info(s)
	}

	if g.HTTPServer != nil {
		g.HTTPServer.Addr = optimiseListenAddress(g.HTTPServer.Addr)
		logger.Info("HTTP Server listening on: ", g.HTTPServer.Addr)
		go g.HTTPServer.ListenAndServe()
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)

	// Register our signal handlers in order
	// ^C in the terminal will send an os.Interrupt
	signal.Notify(quit, os.Interrupt)
	// Kubernetes will send a SIGTERM, so notify on that as well
	signal.Notify(quit, syscall.SIGTERM)

	foundSignal := <-quit
	g.State.Shutdown()

	if foundSignal == syscall.SIGTERM {
		// If we terminate immediately from a SIGTERM, we still may
		// have incoming connections routed to us by k8s. Instead,
		// disable KeepAlives, and sleep to wait for this to propagate.
		if g.HTTPServer != nil {
			g.HTTPServer.SetKeepAlivesEnabled(false)
		}

		logger.Info("SIGTERM received, starting to shut down")
		time.Sleep(15 * time.Second)
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.ShutdownTimeout)
	defer cancel()

	if g.HTTPServer != nil {
		logger.Info("Gracefully shutting down http server with timeout: ", g.ShutdownTimeout)
		time.Sleep(g.ShutdownTimeout)
		if err := g.HTTPServer.Shutdown(ctx); err != nil {
			logger.Fatal("Error shutting down http server: ", err)
		}
	}

	logger.Info("Server exiting...")
	os.Exit(0)
}
