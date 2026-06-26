package service

//logic cho XPATH -> XML Get-Config qua RPC -> ....
import (
	"context"
	"fmt"
	"log/slog"
	model "netconf-test/internal/modelschema"
	"netconf-test/internal/netconfclient"
	"os"
	"strings"
	"time"
)

type ConfigService struct {
	client *netconfclient.NetconfClient
}

func NewEditConfigService(c *netconfclient.NetconfClient) *ConfigService {
	return &ConfigService{
		client: c,
	}
}

func (s *ConfigService) Execute(ctx context.Context, path, namespace, action, value string) error {
	pathSlice := strings.Fields(path)
	dataPayload := s.buildNestedXML(pathSlice, namespace, action, value)

	if action == "edit" {
		if value == "" {
			slog.Error("Invalid configuration. Edit action requires CONFIG_VALUE")
			os.Exit(1)
		}
		return s.hEdit(ctx, dataPayload, path, namespace, value)
	}

	return s.hShow(ctx, dataPayload, path, namespace)
}

func (s *ConfigService) hShow(ctx context.Context, dataPayload, path, namespace string) error {
	showReq := model.GetConfigRequest{
		Source: model.Source{}, //Automatically create tag <running></running>
		Filter: model.Filter{
			Type:  "subtree",
			Inner: dataPayload,
		},
	}

	// rpcPayload := fmt.Sprintf(`
	// 	<get-config>
	// 		<source><running/></source>
	// 		<filter type="subtree">%s</filter>
	// 	</get-config>`, &dataPayload)

	slog.Info("SHOWING (GET-CONFIG) at:", "namespace", namespace, "path", path)
	fmt.Printf("\n=============================================================================================")
	time.Sleep(3 * time.Second)
	//replyXML, err := s.client.ExecRPC(ctx, model.RawXML(rpcPayload))
	replyXML, err := s.client.ExecRPC(ctx, showReq)
	if err != nil {
		slog.Error("Failed sending SHOW-CONFIG RPC: %v", "error", err)
		os.Exit(1)
	}

	slog.Info("Receiving SHOW-CONFIG reply successfully from the Server")
	fmt.Printf("\n===== XML PAYLOAD FROM SERVER =====\n%v\n==================================================\n", replyXML)
	fmt.Printf("")
	return nil
}

func (s *ConfigService) hEdit(ctx context.Context, dataPayload, path, namespace, value string) error {
	//1-Drafting
	editReq := model.EditConfigRequest{
		Target: model.Target{},
		Config: model.Config{
			Inner: dataPayload,
		},
	}

	slog.Info("DRAFTING (EDIT-CONFIG) -> CANDIDATE DATASTORE...", "namespace", namespace, "path", path, "expected_value", value)

	_, err := s.client.ExecRPC(ctx, editReq)
	if err != nil {
		slog.Error("failed sending EDIT-CONFIG RPC: %v", "error", err)
		os.Exit(1)
	}
	fmt.Printf("\n=======================================================================================")
	slog.Info("\n===== EDIT-CONFIG accepted by Candidate datastore=====\n")

	//2-CommitRequest
	time.Sleep(5 * time.Second)
	slog.Info("\n=====Candidate datastore commiting changes=====\n")
	commitReq := model.CommitRequest{}

	_, errC := s.client.ExecRPC(ctx, commitReq)
	if errC != nil {
		slog.Error("failed sending COMMIT RPC: %v", "error", errC)
		os.Exit(1)
	}

	slog.Info("Commit successful. Configuration applied!")
	fmt.Printf("")
	return nil
}

func (s *ConfigService) buildNestedXML(path []string, namespace, action, value string) string {
	if len(path) == 0 {
		return ""
	}

	var openTags, closeTags strings.Builder

	for i, tag := range path {
		if i == 0 && namespace != "" {
			openTags.WriteString(fmt.Sprintf("<%s xmlns=\"%s\">", tag, namespace))
		} else {
			openTags.WriteString(fmt.Sprintf("<%s>", tag))
		}
	}

	for i := len(path) - 1; i >= 0; i-- {
		closeTags.WriteString(fmt.Sprintf("</%s>", path[i]))
	}

	innerValue := ""
	if action == "edit" {
		innerValue = value
	}

	return openTags.String() + innerValue + closeTags.String()
}
