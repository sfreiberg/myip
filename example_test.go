package myip_test

import (
	"fmt"

	"github.com/sfreiberg/myip"
)

func ExampleGetIP() {
	ip, err := myip.GetIP()
	if err != nil {
		fmt.Printf("Error getting IP: %s\n", err)
		return
	}

	fmt.Printf("Your IP address is: %s\n", ip)
}
