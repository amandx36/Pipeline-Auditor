package client

import (
	"context"
	"fmt"
	"log"
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

	log.Printf("Connecting to CI-LogCollector at %s", address)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Printf(
			"Failed to connect to CI-LogCollector at %s: %v",
			address,
			err,
		)
		return nil, err
	}

	log.Printf("Connected to CI-LogCollector at %s", address)

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

func toProtoPipelineEvent(
	event models.PipelineEvent,
) *collectorpb.PipelineEvent {

	if event.Provider == "" &&
		event.Repository == "" &&
		event.PipelineID == "" {
		return nil
	}

	return &collectorpb.PipelineEvent{
		Provider:     event.Provider,
		Repository:   event.Repository,
		PipelineId:   event.PipelineID,
		RunNumber:    uint64(event.RunNumber),
		RunAttempt:   uint32(event.RunAttempt),
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

func (c *Client) CollectLogs(
	event models.PipelineEvent,
) (*collectorpb.CollectLogsResponse, error) {

	if c == nil || c.client == nil {
		return nil, grpc.ErrClientConnClosing
	}

	pipeline := toProtoPipelineEvent(event)

	if pipeline == nil {
		return nil, fmt.Errorf("pipeline event is empty")
	}

	request := &collectorpb.CollectLogsRequest{
		Pipeline: pipeline,
	}

	log.Printf(
		"Sending CollectLogs RPC: provider=%s repository=%s pipeline_id=%s",
		pipeline.GetProvider(),
		pipeline.GetRepository(),
		pipeline.GetPipelineId(),
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	response, err := c.client.CollectLogs(ctx, request)

	if err != nil {
		log.Printf("CollectLogs RPC failed: %v", err)
		return nil, err
	}

	log.Printf(
		"CollectLogs RPC succeeded: accepted=%t collection_id=%s message=%s",
		response.GetAccepted(),
		response.GetCollectionId(),
		response.GetMessage(),
	)

	return response, nil
}

func CollectLogs(
	event models.PipelineEvent,
) (*collectorpb.CollectLogsResponse, error) {

	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	return client.CollectLogs(event)
}
