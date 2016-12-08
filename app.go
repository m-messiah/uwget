package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
)

func encode_size(message []byte) []byte {
	var_size := make([]byte, 2)
	binary.LittleEndian.PutUint16(var_size, uint16(len(message)))
	return var_size
}

func uwsgi_pack(path, query, host, remote_addr string) []byte {
	if path == "" {
		path = "/"
	}
	uwsgi_params := map[string]string{
		"SERVER_PROTOCOL": "HTTP/1.1",
		"REQUEST_METHOD":  "GET",
		"PATH_INFO":       path,
		"REQUEST_URI":     path,
		"QUERY_STRING":    query,
		"SERVER_NAME":     host,
		"HTTP_HOST":       host,
		"REMOTE_ADDR":     remote_addr,
	}
	var params []byte
	for k, v := range uwsgi_params {
		bytes_k, bytes_v := []byte(k), []byte(v)
		params = append(append(params, encode_size(bytes_k)...), bytes_k...)
		params = append(append(params, encode_size(bytes_v)...), bytes_v...)
	}
	return append(append(append([]byte{0}, encode_size(params)...), 0), params...)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] uwsgi://host:port/path\n\nParameters:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	http_host := flag.String("host", "", "HTTP_HOST")
	remote_addr := flag.String("remote", "127.0.0.1", "remote addr")
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
	uwsgi_request := uwsgi_pack(url.Path, url.RawQuery, *http_host, *remote_addr)
	conn, _ := net.Dial("tcp", url.Host)
	conn.Write(uwsgi_request)
	var response bytes.Buffer
	io.Copy(&response, conn)
	fmt.Print(response.String())
	defer conn.Close()
}
