package client

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	collectorpb "github.com/amandx36/Pipeline-Auditor/gen/collector"
	"github.com/amandx36/Pipeline-Auditor/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const collectorAddressEnvKey = "CI_LOG_COLLECTOR_ADDR"

var (
	clientOnce      sync.Once
	sharedClient    *Client
	sharedClientErr error
)

type Client struct {
	conn   *grpc.ClientConn
	client collectorpb.LogCollectorClient
}

func collectorAddress() string {
	addr := strings.TrimSpace(os.Getenv(collectorAddressEnvKey))
	if addr == "" {
		return "localhost:50051"
	}
	return addr
}

func newClient() (*Client, error) {
	address := collectorAddress()
	conn, err := grpc.DialContext(
		context.Background(),
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: collectorpb.NewLogCollectorClient(conn),
	}, nil
}

func GetClient() (*Client, error) {
	clientOnce.Do(func() {
		sharedClient, sharedClientErr = newClient()
	})
	return sharedClient, sharedClientErr
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func toProtoPipelineEvent(event models.PipelineEvent) *collectorpb.PipelineEvent {
	if event.Provider == "" && event.Repository == "" && event.PipelineID == "" {
		return nil
	}

	return &collectorpb.PipelineEvent{
		Provider:     event.Provider,
		Repository:   event.Repository,
		PipelineId:   event.PipelineID,
		RunNumber:    int32(event.RunNumber),
		RunAttempt:   int32(event.RunAttempt),
		PipelineName: event.PipelineName,
		Branch:       event.Branch,
		CommitSha:    event.CommitSHA,
		Event:        event.Event,
		Status:       event.Status,
		Conclusion:   event.Conclusion,
		LogsUrl:      event.LogsURL,
		HtmlUrl:      event.HTMLURL,
	}
}

func (c *Client) CollectLogs(event models.PipelineEvent) (*collectorpb.CollectLogsResponse, error) {
	if c == nil || c.client == nil {
		return nil, grpc.ErrClientConnClosing
	}

	request := &collectorpb.CollectLogsRequest{
		Pipeline: toProtoPipelineEvent(event),
	}
	if request.Pipeline == nil {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.client.CollectLogs(ctx, request)
}

func CollectLogs(event models.PipelineEvent) (*collectorpb.CollectLogsResponse, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	return client.CollectLogs(event)
}
