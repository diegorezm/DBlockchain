package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/diegorezm/DBlockchain/internals/blockchain"
	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
)

type BlockchainServerHandler struct {
	nodes map[string]bool
	mu    sync.Mutex
}

func NewBlockchainServerHandler() *BlockchainServerHandler {
	return &BlockchainServerHandler{
		nodes: make(map[string]bool),
		mu:    sync.Mutex{},
	}
}

// Register a new node inside of the nodes map.
func (s *BlockchainServerHandler) ConnectNode(w http.ResponseWriter, r *http.Request) {
	b, err := webutils.ParseJSON[blockchain.NodeInsert](r.Body)

	if err != nil {
		webutils.WriteBadRequest(w, err.Error())
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nodes[b.Address] = true

	webutils.WriteSuccess[any](w, nil, fmt.Sprintf("Node %s registered successfully.", b.Address))
}

// Register a new node inside of the nodes map.
func (s *BlockchainServerHandler) DisconnectNode(w http.ResponseWriter, r *http.Request) {
	b, err := webutils.ParseJSON[blockchain.NodeInsert](r.Body)

	if err != nil {
		webutils.WriteBadRequest(w, err.Error())
	}

	s.mu.Lock()

	defer s.mu.Unlock()

	delete(s.nodes, b.Address)

	webutils.WriteSuccess[any](w, nil, fmt.Sprintf("Node %s registered successfully.", b.Address))
}

// GetNodes returns the list of all currently registered nodes.
func (s *BlockchainServerHandler) GetNodes(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodeList := make([]string, 0, len(s.nodes))
	for nodeAddr := range s.nodes {
		nodeList = append(nodeList, nodeAddr)
	}

	webutils.WriteSuccess(w, nodeList, "List of registered nodes.")
}

// PingNodes checks the liveness of all registered nodes by sending a GET request to their /ping endpoint.
// If a node does not respond successfully (e.g., network error, timeout, non-200 status),
// it is removed from the list of registered nodes.
func (s *BlockchainServerHandler) PingNodes() {
	log.Println("Starting node ping routine...")

	s.mu.Lock()
	nodesToPing := make([]string, 0, len(s.nodes))
	for nodeAddr := range s.nodes {
		nodesToPing = append(nodesToPing, nodeAddr)
	}
	s.mu.Unlock()
	if len(nodesToPing) == 0 {
		log.Print("No nodes to ping.")
		return
	}

	var wg sync.WaitGroup // WaitGroup to wait for all ping goroutines to complete.
	// Channel to safely collect addresses of nodes that need to be removed.
	nodesToRemove := make(chan string, len(nodesToPing))

	// Iterate over the copied list of nodes and ping each concurrently.
	for _, nodeAddr := range nodesToPing {
		wg.Add(1) // Increment the WaitGroup counter for each goroutine.

		go func(addr string) {
			defer wg.Done() // Decrement the counter when the goroutine finishes.

			pingURL := fmt.Sprintf("%s/ping", addr)
			log.Printf("Pinging node: %s\n", pingURL)

			// Create a context with a timeout for the HTTP request.
			// This prevents pings from hanging indefinitely if a node is truly unresponsive.
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel() // Ensure the context's resources are released.

			req, err := http.NewRequestWithContext(ctx, "GET", pingURL, nil)
			if err != nil {
				log.Printf("Error creating request for %s: %v", addr, err)
				nodesToRemove <- addr // Send node to removal channel
				return
			}

			// Execute the HTTP request using the default client.
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("Node %s did not respond (network error or timeout): %v", addr, err)
				nodesToRemove <- addr // Send node to removal channel
				return
			}

			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				log.Printf("Node %s responded with non-OK status: %d", addr, res.StatusCode)
				nodesToRemove <- addr
				return
			}

			log.Printf("Node %s responded successfully (status %d).", addr, res.StatusCode)

		}(nodeAddr)
	}

	wg.Wait()
	close(nodesToRemove)

	s.mu.Lock()
	defer s.mu.Unlock()

	for addr := range nodesToRemove {
		delete(s.nodes, addr)
		log.Printf("Removed unresponsive node: %s", addr)
	}
	log.Println("Node ping routine finished.")
}

func (s *BlockchainServerHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	webutils.WriteSuccess[any](w, nil, "Pong!")
}
