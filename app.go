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

	http_host := flag.String("host", "", "HTTP_HOST")
	remote_addr := flag.String("remote", "127.0.0.1", "remote addr")
	modifier1 := flag.Int("modifier1", 0, "modifier1")
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
	if *http_host == "" {
		http_host = &host
	}
	response := get(url, *http_host, *remote_addr, *modifier1)
	fmt.Print(string(response))
}
