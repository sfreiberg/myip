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
	ifconfigMe = "http://ifconfig.me"
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

// GetIPUsingIfConfigMe uses the ifconfig.me web service to get your public IP. It returns your public IP or an error.
func GetIPUsingIfConfigMe() (net.IP, error) {
	resp, err := http.Get(ifconfig.me)
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
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read Body")
	}

	ip := net.ParseIP(string(body))
	if ip == nil {
		return nil, fmt.Errorf("Unable to parse IP %q", string(body))
	}
	return ip, nil

}
