package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	gRPC "github.com/DarkLordOfDeadstiny/DSYS-gRPC-template/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type node struct {
	gRPC.UnimplementedTemplateServer
	listenPort string
	name       string
	nodeSlice  []nodeConnection
	mutex      sync.Mutex
}
type nodeConnection struct {
	node     gRPC.TemplateClient
	nodeConn *grpc.ClientConn
}

var nodeName = flag.String("name", "default", "Senders name")
var port = flag.String("port", "5400", "Listen port")

func (n *node) ConnectToNode(port string) {

	//dial options
	//the server is not using TLS, so we use insecure credentials
	//(should be fine for local testing but not in the real world)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	//use context for timeout on the connection
	timeContext, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() //cancel the connection when we are done

	//dial the server to get a connection to it
	log.Printf("client %s: Attempts to dial on port %s\n", n.name, port)
	// Insert your device's IP before the colon in the print statement
	conn, err := grpc.DialContext(timeContext, fmt.Sprintf(":%s", port), opts...)
	if err != nil {
		log.Printf("Fail to Dial : %v", err)
		return
	}

	// makes a client from the server connection and saves the connection
	// and prints rather or not the connection was is READY
	nodeConnection := nodeConnection{
		node:     gRPC.NewTemplateClient(conn),
		nodeConn: conn,
	}
	fmt.Println(nodeConnection)
	n.nodeSlice = append(n.nodeSlice, nodeConnection)
	fmt.Println(n.nodeSlice)
	log.Println("the connection is: ", conn.GetState().String())
}

func (n *node) launchNode() {
	log.Printf("Server %s: Attempts to create listener on port %s\n", n.name, n.listenPort)

	// Create listener tcp on given port or default port 5400
	// Insert your device's IP before the colon in the print statement
	list, err := net.Listen("tcp", fmt.Sprintf(":%s", n.listenPort))
	if err != nil {
		log.Printf("Server %s: Failed to listen on port %s: %v", n.name, n.listenPort, err) //If it fails to listen on the port, run launchServer method again with the next value/port in ports array
		return
	}

	// makes gRPC server using the options
	// you can add options here if you want or remove the options part entirely
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	gRPC.RegisterTemplateServer(grpcServer, n) //Registers the server to the gRPC server.

	log.Printf("Server %s: Listening on port %s\n", n.name, n.listenPort)

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
	// code here is unreachable because grpcServer.Serve occupies the current thread.
}

func (n *node) parseInput() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("--------------------")

	//Infinite loop to listen for clients input.
	for {
		fmt.Print("-> ")

		//Read input into var input and any errors into err
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input) //Trim input
		if strings.Contains(input, "connect") {
			portString := input[8:12]
			if err != nil {
				// ... handle error
				panic(err)
			}
			n.ConnectToNode(portString)
		} else {
			n.sendMessage(input)
		}
		continue
	}
}

func main() {
	//parse flag/arguments
	flag.Parse()

	fmt.Println("--- CLIENT APP ---")

	//log to file instead of console
	//setLog()

	//connect to server and close the connection when program closes
	fmt.Println("--- join Server ---")

	node := node{
		listenPort: *port,
		name:       *nodeName,
		nodeSlice:  make([]nodeConnection, 0),
	}
	go node.launchNode()
	node.parseInput()
}

func (n *node) SendMessage(ctx context.Context, Message *gRPC.Message) (*gRPC.Message, error) {
	// locks the server ensuring no one else can increment the value at the same time.
	// and unlocks the server when the method is done.
	n.mutex.Lock()
	defer n.mutex.Unlock()
	fmt.Println("s.Message: ", Message)
	return &gRPC.Message{Message: "fucking"}, nil
}

func (n *node) sendMessage(input string) {
	Message := &gRPC.Message{
		Message: input,
	}
	for _, element := range n.nodeSlice {
		Message, err := element.node.SendMessage(context.Background(), Message)
		if err != nil {
			log.Printf("Client %s: no response from the server, attempting to reconnect", n.name)
			log.Println(err)
		}

		fmt.Println(Message.Message)
	}
}
