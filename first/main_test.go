package main

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestMain(t *testing.T) {
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		main()
	}()
	switch <-sign {
	case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
		return
	}
}
