package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"golang.org/x/exp/slices"
)

type Request struct {
	Id   int    `json:"id"`
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body any    `json:"body"`
}
type Response struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body any    `json:"body"`
}

const NODE_COUNT = 30

var nextMsgId = 0

type Node struct {
	payloadMutex   sync.Mutex
	nextMsgIdMutex sync.Mutex
	payload        map[string]float64
	nodeId         string
	nodeIds        []string
}

func NewNode(nodeCount float64) (*Node, error) {
	node := &Node{
		payload: make(map[string]float64),
	}
	return node, nil
}

func (node *Node) increment(delta float64) {
	node.replace(node.nodeId, node.payload[node.nodeId]+delta)
}
func (node *Node) localCounter() float64 {
	var localCounter float64 = 0
	for _, v := range node.payload {
		localCounter += v
	}
	return localCounter
}
func (node *Node) replace(nodeId string, value float64) {
	node.payloadMutex.Lock()
	defer node.payloadMutex.Unlock()
	node.payload[nodeId] = value
}
func (node *Node) merge(payload map[string]float64) {
	for nodeId := range node.payload {
		if node.payload[nodeId] < payload[nodeId] {
			fmt.Fprintf(os.Stderr, "Receive value %v greater than local value %v, process to replace", node.payload[nodeId], payload[nodeId])
			node.replace(nodeId, payload[nodeId])
		}
	}
}

func (node *Node) handle(msgType string, req Request) error {
	body := req.Body.(map[string]any)
	switch msgType {
	case "init":
		{
			node.nodeId = body["node_id"].(string)
			raw_nodeIds := body["node_ids"].([]any)
			for _, id := range raw_nodeIds {
				node.nodeIds = append(node.nodeIds, id.(string))
			}

			// Initialize local counter
			for _, nodeId := range node.nodeIds {
				node.payload[nodeId] = 0
			}
			body["type"] = "init_ok"
			node.reply(req, body)
			node.periodic()
			fmt.Fprintf(os.Stderr, "Node %s initialized", node.nodeId)
		}
	case "add":
		{
			delta := body["delta"].(float64)
			node.increment(delta)
			body["type"] = "add_ok"
			delete(body, "delta")
			node.reply(req, body)
			fmt.Fprintf(os.Stderr, "Add %v to node %s\n", delta, node.nodeId)
		}
	case "read":
		{

			body["type"] = "read_ok"
			body["value"] = node.localCounter()
			node.reply(req, body)
			fmt.Fprintf(os.Stderr, "Read result from node %s\n", node.nodeId)
		}
	case "replicate":
		{
			from := req.Src
			body := req.Body.(map[string]any)
			raw_payload := body["value"].(map[string]any)
			payload := make(map[string]float64, 0)
			if len(raw_payload) != 0 {
				for k, v := range raw_payload {
					payload[k] = v.(float64)
				}
			}
			fmt.Fprintf(os.Stderr, "Receive payload %v \n", payload)
			node.merge(payload)
			fmt.Fprintf(os.Stderr, "Done replicating result from node %s to node %s\n", from, node.nodeId)
		}

	}
	return nil
}
func (node *Node) reply(req Request, reqBody map[string]any) {
	body := req.Body.(map[string]any)

	reqBody["in_reply_to"] = body["msg_id"]

	node.send(req.Src, reqBody)
}
func (node *Node) send(dest string, body any) {
	res := Response{
		Src:  node.nodeId,
		Dest: dest,
		Body: body,
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling request: %s for %v", err, res)
		return
	}
	fmt.Fprintln(os.Stdout, string(jsonBytes))
}
func (node *Node) increaseMsgId() {
	node.nextMsgIdMutex.Lock()
	defer node.nextMsgIdMutex.Unlock()
	nextMsgId += 1

}
func (node *Node) reqReplication(other_nid string) error {
	node.increaseMsgId()
	if !slices.Contains(node.nodeIds, other_nid) {
		return fmt.Errorf("error attempting connect to disconnected peer")

	}
	body := make(map[string]any)
	body["type"] = "replicate"
	body["value"] = node.payload
	req := Request{
		Id:   nextMsgId,
		Src:  node.nodeId,
		Dest: other_nid,
		Body: body,
	}

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error when marshaling replication request: %s for %v", err, req)
	}
	fmt.Fprintln(os.Stdout, string(jsonBytes))
	return nil
}

// a node will send replicating req periodically.
func (node *Node) periodic() {
	go func() {
		for {
			for _, other := range node.nodeIds {
				if other == node.nodeId {
					continue
				}
				node.reqReplication(other)
			}
			time.Sleep(3 * time.Second)
		}
	}()
}

//	func (n *Node) periodic() {
//		for _, node := range node.nodeIds {
//			node.merge()
//		}
//	}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	node, _ := NewNode(NODE_COUNT)
	for scanner.Scan() {
		line := scanner.Text()
		msg := Request{}
		err := json.Unmarshal([]byte(line), &msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			continue
		}
		fmt.Fprintf(os.Stderr, "Received \"%v\\n\"\n", msg)

		body := msg.Body.(map[string]any)
		bodyType := body["type"].(string)
		node.handle(bodyType, msg)
	}
}
