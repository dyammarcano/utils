package util

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func WaitToFinish() {
	msg := fmt.Sprintf("received termination signal, shutting down gracefully...")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Block until a termination signal is received
	<-sigCh
	zap.L().Info(msg)
	fmt.Println(msg)
}
