package netconfclient

//Session n Context Management
import (gossh "golang.org/x/crypto/ssh"
	"nemith.io/netconf"
	ncssh "nemith.io/netconf/transport/ssh"

	"context"
	"fmt"
	"log"
	"time"
	model "netconf-test/internal/modelschema"
)

func NewNetconfClient(host, port, user, pass string) (*NetconfClient, error){
	addr := fmt.Sprintf("%s:%s", host, port)

	sshconfig := &gossh.ClientConfig{
		User: user,
		Auth: []gossh.AuthMethod{
			gossh.Password(pass),
		},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout: 10*time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//1. Create Transport (aka Connection)
	transport, err := ncssh.Dial(ctx, "tcp", addr, sshconfig)
	if err != nil {
		log.Fatalf("Failed to connect transportation: %v", err)
	}

	//2. Create Session from Transport - 1Session -> 1Transport
	session, err := netconf.NewSession(transport)
	if err != nil {
		log.Fatalf("Failed to connect transportation: %v", err)
	}

	//3. Return Client Wrapper
	return &NetconfClient{
		session:session,
	}, nil
}

type NetconfClient struct{
	session *netconf.Session
}

//ExecRPC sends raw XML and return reply
func (c *NetconfClient) ExecRPC(ctx context.Context, payload any) (string, error){
	var reply model.RPCReply
	err := c.session.Exec(ctx, payload, &reply)
	if  err != nil{
		return "", err
	}
	return reply.InnerXML, nil
}

//Close connection
func (c *NetconfClient) Close(ctx context.Context) error {
	if c.session != nil {
		return c.session.Close(ctx)
	}
	return nil
}