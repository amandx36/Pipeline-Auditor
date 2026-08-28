package client

import (
	"context"
	"net"
	"testing"

	collectorpb "github.com/amandx36/Pipeline-Auditor/gen/collector"
	"github.com/amandx36/Pipeline-Auditor/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type testCollectorServer struct {
	collectorpb.UnimplementedLogCollectorServer
	request *collectorpb.CollectLogsRequest
}

func (s *testCollectorServer) CollectLogs(_ context.Context, request *collectorpb.CollectLogsRequest) (*collectorpb.CollectLogsResponse, error) {
	s.request = request
	return &collectorpb.CollectLogsResponse{Accepted: true, CollectionId: "collection-1", Message: "queued"}, nil
}

// Integration Test — sends normalized events over gRPC. [Collector Dispatch]
func TestClientCollectLogs(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	collector := &testCollectorServer{}
	collectorpb.RegisterLogCollectorServer(server, collector)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(func() { server.Stop(); _ = listener.Close() })

	conn, err := grpc.NewClient(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	client := &Client{conn: conn, client: collectorpb.NewLogCollectorClient(conn)}
	event := models.PipelineEvent{Provider: "github", Repository: "owner/repository", PipelineID: "42", RunNumber: 7, RunAttempt: 2, PipelineName: "CI", Branch: "main", CommitSHA: "abc", Event: "push", Status: "completed", Conclusion: "failure"}

	response, err := client.CollectLogs(event)
	if err != nil {
		t.Fatalf("CollectLogs() error = %v", err)
	}
	if !response.GetAccepted() || response.GetCollectionId() != "collection-1" || collector.request == nil || collector.request.GetPipeline().GetRepository() != event.Repository {
		t.Fatalf("unexpected gRPC response/request: response=%+v request=%+v", response, collector.request)
	}
}

// Unit Test — rejects an uninitialized gRPC client. [Nil Client]
func TestClientCollectLogsNilClient(t *testing.T) {
	response, err := (*Client)(nil).CollectLogs(models.PipelineEvent{Provider: "github"})
	if response != nil || err == nil {
		t.Fatalf("CollectLogs() = (%+v, %v), want nil response and error", response, err)
	}
}
