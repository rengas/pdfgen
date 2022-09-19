package service

import (
	"os"
	"os/signal"
)

func Wait(signals ...os.Signal) os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, signals...)
	return <-sig
}
