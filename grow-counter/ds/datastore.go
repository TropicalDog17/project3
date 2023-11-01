package ds

import "fmt"

type Message struct {
	Id   int            `json:"id"`
	Src  string         `json:"src"`
	Dest string         `json:"dest"`
	Body map[string]any `json:"body"`
}
type DataStore struct {
	kv map[string]int
}

var nextResponseId = 0

func (d DataStore) handle(message Message) Message {
	nextResponseId += 1
	response := Message{}
	response.Src = message.Dest
	response.Dest = message.Src
	response.Id = nextResponseId
	response.Body = message.Body
	key, _ := message.Body["key"].(string)
	switch message.Body["type"] {
	case "read":
		val, ok := d.kv[key]
		if ok {
			response.Body["type"] = "read_ok"
			response.Body["value"] = val
		} else {
			response.Body["type"] = "error"
		}
	case "write":
		d.kv[key], _ = message.Body["value"].(int)
		message.Body["type"] = "write_ok"
	case "cas":
		// check if key exists
		from := message.Body["from"].(int)
		to := message.Body["to"].(int)
		create_if_not_exists := message.Body["create_if_not_exists"].(bool)
		val, key_exist := d.kv[key]
		if key_exist {
			if val == from {
				d.kv[key] = to
				response.Body["type"] = "cas_ok"
			} else {
				// Value not match
				response.Body["type"] = "error"
			}
		} else {
			if create_if_not_exists {
				d.kv[key] = to
				response.Body["type"] = "cas_ok"
			} else {
				response.Body["type"] = "error"
			}
		}
	}
	return response
}
func main() {
	fmt.Println("Hello world!")
}
