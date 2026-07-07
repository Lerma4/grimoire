package app

import "runtime"

func runtimeOS() string { return runtime.GOOS }
