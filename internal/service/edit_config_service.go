package service

//logic cho XPATH -> XML Get-Config qua RPC -> ....
import (
	"context"
	"fmt"
	"log/slog"
	model "netconf-test/internal/modelschema"
	"netconf-test/internal/netconfclient"
	"strings"
)

type ConfigService struct {
	client *netconfclient.NetconfClient
}

func NewEditConfigService(c *netconfclient.NetconfClient) *ConfigService{
	return &ConfigService{
		client: c,
	}
}

func (s *ConfigService) Execute(ctx context.Context, path, namespace string) error {
	pathSlice := strings.Fields(path)
	dataPayload := s.buildNestedXML(pathSlice, namespace)

	req := model.GetConfigRequest{
		Source: model.Source{},			//Automatically create tag <running></running>
		Filter: model.Filter{
			Type: "subtree",
			Inner: dataPayload,
		},
	}

	// rpcPayload := fmt.Sprintf(`
	// 	<get-config>
	// 		<source><running/></source>
	// 		<filter type="subtree">%s</filter>
	// 	</get-config>`, &dataPayload)
	slog.Info("SHOWING (GET-CONFIG) at:", "namespace", namespace, "path", path)

	//replyXML, err := s.client.ExecRPC(ctx, model.RawXML(rpcPayload))
	replyXML, err := s.client.ExecRPC(ctx, req)
	if err != nil{
		fmt.Print("failed sending RPC: %w", err)
	}

	slog.Info("Receiving reply successfully from the Server")
	fmt.Printf("\n===== XML PAYLOAD FROM SERVER =====\n%v\n==================================================\n", replyXML)
	fmt.Printf("")
	return nil
}

func (s *ConfigService) buildNestedXML(path []string, namespace string) string{
	if len(path) == 0 {
		return ""
	}

	var openTags, closeTags strings.Builder

	for i, tag := range path {
		if i==0 && namespace != "" {
			openTags.WriteString(fmt.Sprintf("<%s xmlns=\"%s\">", tag, namespace))
		} else {
			openTags.WriteString(fmt.Sprintf("<%s>", tag))
		}
	}

	for i := len(path) - 1; i>=0; i--{
		closeTags.WriteString(fmt.Sprintf("</%s>", path[i]))
	}

	return openTags.String() + closeTags.String()
}