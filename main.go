package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"nemith.io/netconf"
)

func main() {
	host := os.Getenv("NETCONF_HOST")
	port := os.Getenv("NETCONF_PORT")
	username := os.Getenv("NETCONF_USERNAME")
	pass := os.Getenv("NETCONF_PASSWORD")

	addr := fmt.Sprintf("%s:%s", host, port)

	session, err := netconf.DialSSH(
		addr,
		netconf.SSHConfigPassword(
			username,
			password,
		),
	)

	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}

	defer session.Close()

	log.Println("connection successfully")

	for {
		time.Sleep(time.Minute)
	}
}