package GuthiEvent

type NodeEvents int

// just commiting as Sanskar's per se
const (
	NodeConnectedEvent NodeEvents = iota // Broadcast this message to all the connected peers
	NodeDisconnectedEvent
)

type Events struct {
	events []NodeEvents
}

// if the node has been connected through listening, broadcast events to all the connected peers directly
// it will enable the network to reach stable synced state quicker

// tf, go
