package main

import (
	"flag"
	"fmt"
	stdlog "log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
)

var (
	Version   string
	BuildTime string
	BuildTag  string

	addr = flag.String("addr", fromEnv("ADDR", ":10000"), "HTTP bind address")

	logger log.Logger
)

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	if Version == "" {
		Version = "0.0"
	}
	if BuildTime == "" {
		BuildTime = "now"
	}
	if BuildTag == "" {
		BuildTag = "dev"
	}

	// Setup logging
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.NewContext(logger).With("time", log.DefaultTimestamp)
	logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	// Redirect stdout to logger
	stdlog.SetFlags(0)
	stdlog.SetOutput(log.NewStdlibAdapter(logger))

	logger.Log(
		"msg", "Starting",
		"addr", *addr,
		"version", Version,
		"buildTime", BuildTime,
		"buildTag", BuildTag,
	)
	defer logger.Log("msg", "Stopped")

	errc := make(chan error, 1)

	// One goroutine to serve the web pages
	go func() {
		mux := http.DefaultServeMux
		mux.HandleFunc("/", Index)
		errc <- http.ListenAndServe(*addr, mux)
	}()

	// Another goroutine waits for signals
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		errc <- nil
	}()

	if err := <-errc; err != nil {
		logger.Log("err", err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	logger.Log("url", r.URL.String(), "remoteAddr", r.RemoteAddr, "ua", r.UserAgent())

	fmt.Fprintf(w, `
<html>
	<head>
		<meta charset=utf-8>
		<title>Docker Test Web</title>
	</head>
	<body>
		<h1>Docker Test Web %s</h1>
		<p>Build time: %s</p>
		<p>Build tag: %s</p>
	</body>
</html>
`, Version, BuildTime, BuildTag)
}

func fromEnv(name, defvalue string) string {
	if s := os.Getenv(name); s != "" {
		return s
	}
	return defvalue
}
