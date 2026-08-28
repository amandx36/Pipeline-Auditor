package client

import (
	"context"
	"log"
	"time"

	collectorpb "github.com/amandx36/Pipeline-Auditor/gen/collector"
	"github.com/amandx36/Pipeline-Auditor/internal/models"
	"google.golang.org/grpc"
)

func (c *Client) CollectLogs(event models.PipelineEvent) (*collectorpb.CollectLogsResponse, error) {
	// is client
	log.Printf("Client logs exist ")
	if c == nil || c.client == nil {
		log.Printf("not client exist %v", grpc.ErrClientConnClosing)
		return nil, grpc.ErrClientConnClosing
	}
	// convert  the pipeline event into the protobuf model
	pipeline := toProtoPipelineEvent(event)
	if pipeline == nil {
		log.Printf("Pipe line event is empty ")
	}

	// creating the grpc request
	request := &collectorpb.CollectLogsRequest{
		Pipeline: pipeline,
	}
	log.Printf("grpc request created \n Sending CollectLogs RPC: provider=%s repository=%s pipeline_id=%s",
		pipeline.GetProvider(),
		pipeline.GetRepository(),
		pipeline.GetPipelineId())

	// max wait 10 sec no response then cancel it
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()
	// rpc , calls the CollectLogs() method
	// implemented by CI-LogCollector
	response, err := c.client.CollectLogs(
		ctx,
		request,
	)
	if err != nil {

		log.Printf(
			"CollectLogs RPC failed: %v",
			err,
		)

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

// making a  wrapper so calling the collectlog will be easy
func CollectLogWrapper(event models.PipelineEvent) (*collectorpb.CollectLogsResponse, error) {
	// shared grpc client
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	return client.CollectLogs(event)
}
