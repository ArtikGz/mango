package main

import (
	"flag"
	"fmt"
	mango "mango/src/network/tcp"
	"net/http"
	_ "net/http/pprof"

	"mango/src/config"
	"mango/src/logger"
)

var (
	configPath   = flag.String("c", "config.json", "-c /path/to/config_file.json")
	profilerPort = flag.Int("p", 0, "-p <profiler port>")
)

func main() {
	flag.Parse()
	err := config.Parse(*configPath)
	if err != nil {
		logger.Error("Couldn't read config due to an error, using the default configuration... (err: %s)", err)
		config.LoadDefaultConfig()
	}

	// profiler
	if *profilerPort != 0 {
		go runProfiler()
	}

	server, err := mango.NewTcpServer(config.Host(), config.Port())
	if err != nil {
		logger.Fatal("Server couldn't start")
	}

	server.Start()
}

func runProfiler() {
	logger.Info("Profiling server runing on 127.0.0.1:%d...", *profilerPort)
	logger.Info(http.ListenAndServe(fmt.Sprintf(":%d", *profilerPort), nil))
}
