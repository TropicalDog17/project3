package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func get_initial_ds() DataStore {
	ds := DataStore{}
	ds.kv = map[string]int{
		"a": 12,
		"b": 34,
		"c": 56,
		"d": 78,
		"e": 910,
	}
	return ds
}
func TestHandleReadData(t *testing.T) {
	ds := get_initial_ds()
	// Id   int            `json:"id"`
	// Src  string         `json:"src"`
	// Dest string         `json:"dest"`
	// Body map[string]any `json:"body"`
	msg := Message{
		Id:   1,
		Src:  "n0",
		Dest: "n2",
		Body: map[string]any{
			"type": "read",
			"key":  "a",
		},
	}
	response := ds.handle(msg)
	assert.Equal(t, response.Body["type"], "read_ok", "Should be ok")
	assert.Equal(t, response.Body["value"], 12, "Should pass")

}
func TestHandleWriteData(t *testing.T) {
	ds := get_initial_ds()
	msg := Message{
		Id:   1,
		Src:  "n0",
		Dest: "n2",
		Body: map[string]any{
			"type":  "write",
			"key":   "a",
			"value": 34,
		},
	}
	response := ds.handle(msg)
	assert.Equal(t, response.Body["type"], "write_ok", "Should be ok")
	assert.Equal(t, response.Body["value"], 34, "Should pass")
	assert.Equal(t, ds.kv["a"], 34, "should equals")
	// Non existing key
	msg = Message{
		Id:   1,
		Src:  "n0",
		Dest: "n2",
		Body: map[string]any{
			"type":  "write",
			"key":   "f",
			"value": 34,
		},
	}
	response = ds.handle(msg)
	assert.Equal(t, response.Body["type"], "write_ok", "Should be ok")
	assert.Equal(t, response.Body["value"], 34, "Should pass")
	assert.Equal(t, ds.kv["f"], 34, "should equals")
}
func TestHandleCas(t *testing.T) {
	ds := get_initial_ds()
	msg := Message{
		Id:   1,
		Src:  "n0",
		Dest: "n2",
		Body: map[string]any{
			"type":                 "cas",
			"key":                  "a",
			"from":                 12,
			"to":                   34,
			"create_if_not_exists": false,
		},
	}
	response := ds.handle(msg)
	assert.Equal(t, response.Body["type"], "cas_ok", "Should be ok")
	assert.Equal(t, ds.kv["a"], 34, "should equals")

	// Test mismatch from value
	ds = get_initial_ds()
	msg = Message{
		Id:   1,
		Src:  "n0",
		Dest: "n2",
		Body: map[string]any{
			"type":                 "cas",
			"key":                  "a",
			"from":                 13,
			"to":                   34,
			"create_if_not_exists": false,
		},
	}
	response = ds.handle(msg)
	assert.Equal(t, response.Body["type"], "error")

	// Test key not exist, create_if_not_exists = false
	ds = get_initial_ds()
	msg = Message{
		Id:   1,
		Src:  "n0",
		Dest: "n2",
		Body: map[string]any{
			"type":                 "cas",
			"key":                  "dsl",
			"from":                 13,
			"to":                   34,
			"create_if_not_exists": false,
		},
	}
	response = ds.handle(msg)
	assert.Equal(t, response.Body["type"], "error")
	// Test key not exist, create_if_not_exists = true
	ds = get_initial_ds()
	msg = Message{
		Id:   1,
		Src:  "n0",
		Dest: "n2",
		Body: map[string]any{
			"type":                 "cas",
			"key":                  "dsl",
			"from":                 13,
			"to":                   34,
			"create_if_not_exists": true,
		},
	}
	response = ds.handle(msg)
	assert.Equal(t, response.Body["type"], "cas_ok", "Should be ok")
	assert.Equal(t, ds.kv["dsl"], 34, "should equals")

}
