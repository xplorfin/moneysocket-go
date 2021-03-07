package message

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/request"
)

// keys
const MessageClass = base.MessageClassKey
const NotificationName = "notification_name"
const RequestName = request.RequestNameKey

// notification types
const (
	NotifyProvider         = "NOTIFY_PROVIDER"
	NotifyProviderNotReady = "NOTIFY_PROVIDER_NOT_READY"
	NotifyPing             = "NOTIFY_PING"
	NotifyPong             = "NOTIFY_PONG"
	// rendezcous notificaton types
	NotifyRendezvous         = "NOTIFY_RENDEZVOUS"
	NotifyRendezvousNotReady = "NOTIFY_RENDEZVOUS_NOT_READY"
	NotifyRendezvousEnd      = "NOTIFY_RENDEZVOUS_END"
)

// nexus lifecycle events
const (
	NexusCreated   = "NEXUS_CREATED"
	NexusAnnounced = "NEXUS_ANNOUNCED"
	NexusWaiting   = "NEXUS_WAITING"
	NexusDestroyed = "NEXUS_DESTROYED"
	NexusRevoked   = "NEXUS_REVOKED"
)

// layers

const (
	Consumer           = "CONSUMER"
	Relay              = "RELAY"
	OutgoingLocal      = "OUTGOING_LOCAL"
	IncomingWebsocket  = "INCOMING_WEBSOCKET"
	IncomingLocal      = "INCOMING_LOCAL"
	OutgoingWebsocket  = "OUTGOING_WEBSOCKET"
	OutgoingRendezvous = "OUTGOING_RENDEZVOUS"
	IncomingRendezvous = "INCOMING_RENDEZVOUS"
	RequestRendezvous  = "REQUEST_RENDEZVOUS"
	Provider           = "PROVIDER"
	ProviderTransact   = "PROVIDER_TRANSACT"
	Terminus           = "TERMINUS"
	Seller             = "SELLER"
)
