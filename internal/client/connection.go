package client 
// only job is connect and create and reuse the connection 
import (
		"context"
		"log"
		"os"
		"sync"
		"time"

		"google.golang.org/grpc"
		"google.golang.org/grpc/credentials/insecure"  // grpc connection without encryption 

	collectorpb "github.com/amandx36/Pipeline-Auditor/gen/collector"
)
// need the collection address 
var collectorAddress = os.Getenv("CI_LOG_COLLECTOR_ADDR")

type Client struct{
	// actual grpc connection 
	conn *grpc.ClientConn ;
	// client is the generated gRPC client.
	// It provides methods to call the LogCollector service
	client collectorpb.LogCollectorClient;

}
func init() {
    	if collectorAddress == "" {
        	collectorAddress = "localhost:50501"
    	}
	}

var (
	clientMu sync.Mutex;
	sharedClient *Client ;
)

// create new connection to  grpc ci-log-collector 
func newClient()(*Client, error){
	var addres = collectorAddress ;
	log.Println("Start the connection to ci-log-collector 01 %s",addres);
	// creating a context of 5 sec 
	ctx,cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel();
	// create the grpc connection 
	connection,err := grpc.DialContext(
		ctx,
		addres,
		// not using the tls  right now 
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		// wait until the connectio get ready 
		grpc.WithBlock(),
	)
	if err != nil {
		log.Println("Error while connecting to the ci - log collector ")
		return nil , err 
	}
	log.Printf(
		"[PIPELINE-AUDITOR] Connected to CI-LogCollector at %s",
		addres,
	)
	// create the client 
	return &Client{
		conn:   connection,
		client: collectorpb.NewLogCollectorClient(connection),
	}, nil

}

// if  already has connection return it other  wise create one 
func GetClient()(*Client, error){
	// lock it  other wise multiple goroutine can create it 
	clientMu.Lock()
	// when function finish unlock it 
	defer clientMu.Unlock()

	// Do we already have a client?
	if sharedClient != nil && sharedClient.conn != nil {
		return sharedClient, nil
	}

	// No client exists.
	//
	// Create a new connection.
	client, err := newClient()

	if err != nil {
		// Don't store failed connections.
		return nil, err
	}

	// Store the  client
	sharedClient = client

	return sharedClient, nil
}
// close  the grpc when no it use 
func (c *Client)Close() error{
	if c == nil || c.conn == nil {
		return nil
	}

	return c.conn.Close()
}