package myip

import (
	"testing"
)

func TestGetIP(t *testing.T) {
	_, err := GetIP()
	if err != nil {
		t.Errorf("Unable to get ip: %s", err)
	}
}
