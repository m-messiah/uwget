package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] uwsgi://host:port/path\n\nParameters:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	httpHost := flag.String("host", "", "HTTP_HOST")
	remoteAddr := flag.String("remote", "127.0.0.1", "remote addr")
	modifier1 := flag.Int("modifier1", 0, "modifier1")
	expectedStatus := flag.String("expected-status", "", "Fail if response status not equal")
	quiet := flag.Bool("q", false, "Disable output")
	flag.Parse()
	arg := flag.Arg(0)
	if arg == "" {
		flag.Usage()
	}
	url, err := url.Parse(arg)
	if err != nil {
		flag.Usage()
	}
	host, _, _ := net.SplitHostPort(url.Host)
	if *httpHost == "" {
		httpHost = &host
	}
	response := get(url, *httpHost, *remoteAddr, *modifier1)
	if *expectedStatus != "" {
		os.Exit(int(checkStatus(response, *expectedStatus)))
	}
	if !*quiet {
		fmt.Print(string(response))
	}
}
