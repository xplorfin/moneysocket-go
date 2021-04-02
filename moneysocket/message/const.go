package message

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// MessageClass is the class of the message.
const MessageClass = base.MessageClassKey

// nexus lifecycle events.
const (
	// NexusCreated announces a nexus has been created.
	NexusCreated = "NEXUS_CREATED"
	// NexusAnnounced announces a nexus has been announced.
	NexusAnnounced = "NEXUS_ANNOUNCED"
	// NexusWaiting announces a nexus has been announced.
	NexusWaiting = "NEXUS_WAITING"
	// NexusDestroyed announces a nexus has been destroyed.
	NexusDestroyed = "NEXUS_DESTROYED"
	// NexusDestroyed announces a nexus has been destroyed.
	NexusRevoked = "NEXUS_REVOKED"
)

// layers

const (
	// Consumer layer name.
	Consumer = "CONSUMER"
	// Relay layer name.
	Relay = "RELAY"
	// OutgoingLocal layer name.
	OutgoingLocal = "OUTGOING_LOCAL"
	// IncomingWebsocket layer name.
	IncomingWebsocket = "INCOMING_WEBSOCKET"
	// IncomingLocal layer name.
	IncomingLocal = "INCOMING_LOCAL"
	// OutgoingWebsocket layer name.
	OutgoingWebsocket = "OUTGOING_WEBSOCKET"
	// OutgoingRendezvous layer name.
	OutgoingRendezvous = "OUTGOING_RENDEZVOUS"
	// IncomingRendezvous layer name.
	IncomingRendezvous = "INCOMING_RENDEZVOUS"
	// RequestRendezvous  layer name.
	RequestRendezvous = "REQUEST_RENDEZVOUS"
	// Provider layer name.
	Provider = "PROVIDER"
	// ProviderTransact layer name.
	ProviderTransact = "PROVIDER_TRANSACT"
	// Terminus layer name.
	Terminus = "TERMINUS"
	// Seller layer name.
	Seller = "SELLER"
)
