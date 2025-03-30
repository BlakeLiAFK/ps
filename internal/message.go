package internal

import "encoding/json"

type Message struct {
	ID        int64
	Data      json.RawMessage
	Timestamp int64
}
