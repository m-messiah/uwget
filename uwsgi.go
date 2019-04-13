package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"net/url"
	"regexp"
)

func encodeSize(message []byte) []byte {
	size := make([]byte, 2)
	binary.LittleEndian.PutUint16(size, uint16(len(message)))
	return size
}

func uwsgiPack(path, query, host, remoteAddr string, modifier1 int) []byte {
	if path == "" {
		path = "/"
	}
	uwsgiParams := map[string]string{
		"SERVER_PROTOCOL": "HTTP/1.1",
		"REQUEST_METHOD":  "GET",
		"PATH_INFO":       path,
		"REQUEST_URI":     path,
		"QUERY_STRING":    query,
		"SERVER_NAME":     host,
		"HTTP_HOST":       host,
		"REMOTE_ADDR":     remoteAddr,
	}
	var params []byte
	for k, v := range uwsgiParams {
		bytesKey, bytesValue := []byte(k), []byte(v)
		params = append(append(params, encodeSize(bytesKey)...), bytesKey...)
		params = append(append(params, encodeSize(bytesValue)...), bytesValue...)
	}
	return append(append(append([]byte{byte(modifier1)}, encodeSize(params)...), 0), params...)
}

func get(url *url.URL, httpHost, remoteAddr string, modifier1 int) []byte {
	uwsgiRequest := uwsgiPack(url.Path, url.RawQuery, httpHost, remoteAddr, modifier1)
	conn, err := net.Dial("tcp", url.Host)
	if err != nil {
		return []byte("No connection")
	}
	conn.Write(uwsgiRequest)
	var response bytes.Buffer
	io.Copy(&response, conn)
	defer conn.Close()
	return response.Bytes()
}

func checkStatus(response []byte, expectedStatus string) int {
	responseRe := regexp.MustCompile(`^HTTP/[01].[01] ([0-9]{3}) `)
	matches := responseRe.FindSubmatch(response)
	if len(matches) < 2 {
		return 1
	}
	if string(matches[1]) == expectedStatus {
		return 0
	}
	return 1
}
