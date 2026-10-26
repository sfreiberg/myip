// Package myip is a package for getting your public IP address.
package myip

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// Service endpoints. These are variables so tests can point them at a local
// server.
var (
	cloudflare = "https://1.1.1.1/cdn-cgi/trace"
	ifconfigMe = "https://ifconfig.me/ip"
)

// client is used for all requests so a stalled service cannot block the
// caller forever.
var client = &http.Client{Timeout: 15 * time.Second}

// GetIP returns your public IP or an error
func GetIP() (net.IP, error) {
	resp, err := client.Get(cloudflare)
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

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read the response body: %s", err)
	}

	lines := bytes.SplitSeq(out, []byte("\n"))
	for line := range lines {
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

// GetIPUsingIfConfigMe uses the ifconfig.me web service to get your public IP. It returns your public IP or an error.
func GetIPUsingIfConfigMe() (net.IP, error) {
	resp, err := client.Get(ifconfigMe)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf(
			"Received unexpected status code from ifconfigMe, expected %v but got %v",
			200,
			resp.StatusCode,
		)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read the response body: %s", err)
	}

	ipS := string(bytes.TrimSpace(body))
	ip := net.ParseIP(ipS)
	if ip == nil {
		return nil, fmt.Errorf("Unable to parse IP %q", ipS)
	}
	return ip, nil

}
