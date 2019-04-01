// Package myip is a package for getting your public IP address.
package myip

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

const (
	cloudflare = "http://1.1.1.1/cdn-cgi/trace"
)

// GetIP returns your public IP or an error
func GetIP() (net.IP, error) {
	resp, err := http.Get(cloudflare)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf(
			"Received unexpected status code from cloudflare, expected %v but got %v",
			200,
			resp.StatusCode,
		)
		return nil, err
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read the response body: %s", err)
	}

	lines := bytes.Split(out, []byte("\n"))
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte("ip=")) {
			ipS := string(line[3:])
			ip := net.ParseIP(ipS)
			if ip == nil {
				return nil, fmt.Errorf("Unable to parse IP %q", ipS)
			}
			return ip, nil
		}
	}

	return nil, fmt.Errorf("Unable to find IP")
}
