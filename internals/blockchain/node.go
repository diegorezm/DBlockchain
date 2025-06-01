package blockchain

import "net/url"

type Node struct {
	Address *url.URL `json:"address"`
}

type NodeBulkRequest struct {
	Nodes []string `json:"nodes"`
}

func NewNode(addr string) (*Node, error) {
	address, err := url.Parse(addr)

	if err != nil {
		return nil, err
	}

	return &Node{Address: address}, nil
}
