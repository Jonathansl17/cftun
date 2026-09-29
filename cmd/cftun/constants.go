package main

import (
	"os"
	"syscall"
)

var version = "dev"

var shutdownSignals = []os.Signal{syscall.SIGTERM}
