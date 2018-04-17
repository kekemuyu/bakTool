package main

import (
	"bakTool/sync"
	"fmt"
	"runtime"
)

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func main() {
	ConfigRuntime()
	src := "tmp1"
	des := []string{"tmp2", "tmp3"}
	sync := sync.New(src, des)
	sync.Run()
}
