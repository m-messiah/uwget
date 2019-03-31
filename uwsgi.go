package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"net/url"
	"regexp"
)

func encode_size(message []byte) []byte {
	var_size := make([]byte, 2)
	binary.LittleEndian.PutUint16(var_size, uint16(len(message)))
	return var_size
}

func uwsgi_pack(path, query, host, remote_addr string, modifier1 int) []byte {
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
	return append(append(append([]byte{byte(modifier1)}, encode_size(params)...), 0), params...)
}

func get(url *url.URL, http_host, remote_addr string, modifier1 int) []byte {
	uwsgi_request := uwsgi_pack(url.Path, url.RawQuery, http_host, remote_addr, modifier1)
	conn, err := net.Dial("tcp", url.Host)
	if err != nil {
		return []byte("No connection")
	}
	conn.Write(uwsgi_request)
	var response bytes.Buffer
	io.Copy(&response, conn)
	defer conn.Close()
	return response.Bytes()
}

func check_status(response []byte, expected_status string) int {
	response_re := regexp.MustCompile(`^HTTP/[01].[01] ([0-9]{3}) `)
	matches := response_re.FindSubmatch(response)
	if len(matches) < 2 {
		return 1
	}
	if string(matches[1]) == expected_status {
		return 0
	}
	return 1
}
