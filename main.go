package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"context"

	gossh "golang.org/x/crypto/ssh"
	"nemith.io/netconf"
	ncssh "nemith.io/netconf/transport/ssh"
)

func main() {
	host := os.Getenv("NETCONF_HOST")
	port := os.Getenv("NETCONF_PORT")
	username := os.Getenv("NETCONF_USERNAME")
	password := os.Getenv("NETCONF_PASSWORD")

	addr := fmt.Sprintf("%s:%s", host, port)

	config := &gossh.ClientConfig{
		User: username,
		Auth: []gossh.AuthMethod{
			gossh.Password(password),
		},
		Timeout: 10 * time.Second,
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()
	
	transport, err := ncssh.Dial(
		ctx,
		"tcp",
		addr,
		config,
	)
	if err != nil {
		log.Fatalf("failed to connect transport: %v", err)
	}
	defer transport.Close()

	session, err := netconf.NewSession(
		transport,
	)
	if err != nil {
		log.Fatalf("failed to create NETCONF session: %v", err)
	}
	defer session.Close(context.Background())

	log.Println("connection successfully")

	// giữ pod sống để xem log
	for {
		time.Sleep(time.Minute)
	}
}