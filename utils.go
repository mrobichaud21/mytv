package main

import (
	"fmt"
	"net"

	"github.com/spf13/viper"
)

func getTCPAddr(key string) *net.TCPAddr {
	addr, addrErr := net.ResolveTCPAddr("tcp", viper.GetString(key))
	if addrErr != nil {
		panic(fmt.Errorf("error parsing address %s: %s", viper.GetString(key), addrErr))
	}
	return addr
}
