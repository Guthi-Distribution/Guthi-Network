package nodes

// var net_platform NetworkPlatform

// message for get address request
// request for all the known address
// maybe we could handle nodes directly???
type GetAddress struct {
	AddrFrom   string
	message_id uint64
}

type GetNodes struct {
	AddrFrom string
	Address  []string
}

// send node object message
type NodesMessage struct {
	AddrFrom string
	Nodes    []NetworkNode // array to make is generic
}

type RequestMessage struct {
	AddrFrom string
}

type AckMessage struct {
	AddrFrom string
}
