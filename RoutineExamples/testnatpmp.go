package main

import (
	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
	"fmt"
)

func main() {
	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		return
	}

	fmt.Print("got ip", gatewayIP.String())
	client := natpmp.NewClient(gatewayIP)
	response, err := client.GetExternalAddress()
	if err != nil {
		fmt.Errorf("error info %s\n", err.Error())
		return
	}
	fmt.Print("External IP address:", response.ExternalIPAddress)
}
