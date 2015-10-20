package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var address string

func init() {
	flag.StringVar(&address, "remote", "127.0.0.1:2000", "Remote address")
}

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	msgData := make([]byte, 2048, 2048)

	for {
		i, err := conn.Read(msgData)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Printf("%d % x\n", time.Now().UnixNano(), msgData[0:i])
	}
}
