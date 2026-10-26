package myip

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// serve points *endpoint at a test server that replies with status and body
// for the duration of the test.
func serve(t *testing.T, endpoint *string, status int, body string) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	orig := *endpoint
	*endpoint = srv.URL
	t.Cleanup(func() { *endpoint = orig })
}

func TestGetIP(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		want    string
		wantErr bool
	}{
		{name: "ipv4", status: 200, body: "fl=1\nip=203.0.113.7\nts=1\n", want: "203.0.113.7"},
		{name: "ipv6", status: 200, body: "fl=1\nip=2001:db8::1\nts=1\n", want: "2001:db8::1"},
		{name: "full trace", status: 200, body: "fl=29f45\nh=1.1.1.1\nip=198.51.100.23\nts=1700000000.123\nvisit_scheme=http\ncolo=SJC\n", want: "198.51.100.23"},
		{name: "bad status", status: 503, body: "ip=203.0.113.7\n", wantErr: true},
		{name: "no ip line", status: 200, body: "fl=1\nts=1\n", wantErr: true},
		{name: "unparsable ip", status: 200, body: "ip=not-an-ip\n", wantErr: true},
		{name: "empty body", status: 200, body: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serve(t, &cloudflare, tt.status, tt.body)
			ip, err := GetIP()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got ip %v", ip)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !ip.Equal(net.ParseIP(tt.want)) {
				t.Errorf("got %v; want %v", ip, tt.want)
			}
		})
	}
}

func TestGetIPUsingIfConfigMe(t *testing.T) {
	serve(t, &ifconfigMe, 200, "192.0.2.44\n")
	ip, err := GetIPUsingIfConfigMe()
	if err != nil {
		t.Fatal(err)
	}
	if !ip.Equal(net.ParseIP("192.0.2.44")) {
		t.Errorf("got %v", ip)
	}
}

func TestGetIPUsingIfConfigMeErrors(t *testing.T) {
	t.Run("bad status", func(t *testing.T) {
		serve(t, &ifconfigMe, 500, "192.0.2.44")
		if _, err := GetIPUsingIfConfigMe(); err == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("html body", func(t *testing.T) {
		serve(t, &ifconfigMe, 200, "<!DOCTYPE html><html></html>")
		if _, err := GetIPUsingIfConfigMe(); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestEndpointsUseHTTPS(t *testing.T) {
	for _, u := range []string{cloudflare, ifconfigMe} {
		if parsed, err := url.Parse(u); err != nil || parsed.Scheme != "https" {
			t.Errorf("%s: want an https endpoint so the answer cannot be spoofed in transit", u)
		}
	}
}

// TestLiveGetIP queries the real Cloudflare endpoint. It is skipped with
// -short or when the network is unreachable.
func TestLiveGetIP(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live network test in -short mode")
	}
	ip, err := GetIP()
	if _, ok := errors.AsType[*url.Error](err); ok {
		t.Skipf("network unavailable: %v", err)
	}
	if err != nil {
		t.Fatalf("Unable to get ip: %s", err)
	}
	if ip == nil {
		t.Fatal("got nil IP")
	}
}

// TestLiveGetIPUsingIfConfigMe queries the real ifconfig.me endpoint. It is
// skipped with -short or when the network is unreachable.
func TestLiveGetIPUsingIfConfigMe(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live network test in -short mode")
	}
	ip, err := GetIPUsingIfConfigMe()
	if _, ok := errors.AsType[*url.Error](err); ok {
		t.Skipf("network unavailable: %v", err)
	}
	if err != nil {
		t.Fatalf("Unable to get ip: %s", err)
	}
	if ip == nil {
		t.Fatal("got nil IP")
	}
}
