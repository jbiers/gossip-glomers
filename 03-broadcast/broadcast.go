package main

const (
	BroadcastType   = "broadcast"
	BroadcastOKType = "broadcast_ok"
	ReadType        = "read"
	ReadOKType      = "read_ok"
	TopologyType    = "topology"
	TopologyOKType  = "topology_ok"
)

type BroadcastBody struct {
	Type      string   `json:"type,omitempty"`
	MsgID     int      `json:"msg_id,omitempty"`
	InReplyTo int      `json:"in_reply_to,omitempty"`
	Code      int      `json:"code,omitempty"`
	Text      string   `json:"text,omitempty"`
	Message   *float64 `json:"message,omitempty"`
}

type ReadBody struct {
	Type      string     `json:"type,omitempty"`
	MsgID     int        `json:"msg_id,omitempty"`
	InReplyTo int        `json:"in_reply_to,omitempty"`
	Code      int        `json:"code,omitempty"`
	Text      string     `json:"text,omitempty"`
	Messages  *[]float64 `json:"messages,omitempty"`
}

type TopologyBody struct {
	Type      string              `json:"type,omitempty"`
	MsgID     int                 `json:"msg_id,omitempty"`
	InReplyTo int                 `json:"in_reply_to,omitempty"`
	Code      int                 `json:"code,omitempty"`
	Text      string              `json:"text,omitempty"`
	Topology  map[string][]string `json:"topology,omitempty"`
}
