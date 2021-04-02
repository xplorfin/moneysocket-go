package layer

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

// OnLayerEventFn is a function downstream layers call on an event.
type OnLayerEventFn = func(layerName string, nexus nexusHelper.Nexus, event string)

// OnAnnounceFn is the announce function for the nexus.
type OnAnnounceFn = func(nexus nexusHelper.Nexus)

// OnRevokeFn is the revoke function for the nexus.
type OnRevokeFn = func(nexus nexusHelper.Nexus)

// Base is the base class for all layers.
type Base interface {
	// SetOnLayerEvent set on layer event
	SetOnLayerEvent(o OnLayerEventFn)
	// SetOnAnnounce on announce event
	SetOnAnnounce(o OnAnnounceFn)
	// SetOnRevoke on revoke event
	SetOnRevoke(o OnAnnounceFn)
	// RegisterAboveLayer above layer events with current layer
	// must be done here since announce nexus
	// is not available form base layer
	RegisterAboveLayer(belowLayer Base)
	// RegisterLayerEvent layer event with nexuses
	RegisterLayerEvent(fn OnLayerEventFn, layerName string)
	// AnnounceNexus nexusHelper
	AnnounceNexus(belowNexus nexusHelper.Nexus)
	// TrackNexus nexusHelper events
	TrackNexus(nexus nexusHelper.Nexus, belowNexus nexusHelper.Nexus)
	// UntrackNexus remove's tracker from nexusHelper events
	UntrackNexus(nexus nexusHelper.Nexus, belowNexus nexusHelper.Nexus)
	// TrackNexusAnnounced tracks that a nexus has been announced
	TrackNexusAnnounced(nexus nexusHelper.Nexus)
	// IsNexusAnnounced checks wether or not a nexus has been announced
	IsNexusAnnounced(nexus nexusHelper.Nexus) bool
	// RevokeNexus revokes a nexus
	RevokeNexus(belowNexus nexusHelper.Nexus)
}

// BaseLayer is used as a superclass for layers.
type BaseLayer struct {
	// LayerName is the name of the current layer this is
	// a string rather than a method call to make debugging easier
	LayerName string
	// OnLayerEvent is a nullable function to be called when a layer
	// event occurs
	OnLayerEvent OnLayerEventFn
	// OnAnnounce is called when a nexus is announced to the layer (from below)
	OnAnnounce OnAnnounceFn
	// OnRevoke is called when a nexus is revoked from the layer (from below)
	OnRevoke OnRevokeFn
	// Nexuses is a thread-safe map of Nexuses to their ids uuid[nexus]
	Nexuses NexusMap
	// Announced is a thread-safe map of Nexuses to their ids uuid[nexus]
	Announced NexusMap
	// BelowNexuses is a thread-safe map of BelowNexuses to their ids uuid[nexus]
	BelowNexuses NexusMap
	// NexusByBelow is a map of below nexuses by nexus[uuid]
	NexusByBelow NexusUUIDMap
	// BelowByNexus is a map of nexus->uuid
	BelowByNexus NexusUUIDMap
}

// SetOnLayerEvent sets the layer event.
func (l *BaseLayer) SetOnLayerEvent(o OnLayerEventFn) {
	l.OnLayerEvent = o
}

// SetOnAnnounce sets an on announce function.
func (l *BaseLayer) SetOnAnnounce(o OnAnnounceFn) {
	l.OnAnnounce = o
}

// SetOnRevoke sets the revoke callback.
func (l *BaseLayer) SetOnRevoke(o OnRevokeFn) {
	l.OnRevoke = o
}

// RegisterLayerEvent registers the layer event.
func (l *BaseLayer) RegisterLayerEvent(fn OnLayerEventFn, layerName string) {
	l.LayerName = layerName
	l.OnLayerEvent = fn
}

// TrackNexus trackers the nexus.
func (l *BaseLayer) TrackNexus(nexus nexusHelper.Nexus, belowNexus nexusHelper.Nexus) {
	l.Nexuses.Store(nexus.UUID(), nexus)
	l.BelowNexuses.Store(belowNexus.UUID(), belowNexus)
	l.NexusByBelow.Store(belowNexus.UUID(), nexus.UUID())
	l.BelowByNexus.Store(nexus.UUID(), belowNexus.UUID())
	l.SendLayerEvent(nexus, message.NexusCreated)
}

// UntrackNexus untracks the nexus.
func (l *BaseLayer) UntrackNexus(nexus nexusHelper.Nexus, belowNexus nexusHelper.Nexus) {
	l.Nexuses.Delete(nexus.UUID())
	l.BelowByNexus.Delete(belowNexus.UUID())
	l.NexusByBelow.Delete(belowNexus.UUID())
	l.BelowByNexus.Delete(belowNexus.UUID())
	l.SendLayerEvent(nexus, message.NexusDestroyed)
}

// TrackNexusAnnounced tracks the nexus.
func (l *BaseLayer) TrackNexusAnnounced(nexus nexusHelper.Nexus) {
	l.Announced.Store(nexus.UUID(), nexus)
}

// IsNexusAnnounced determines wether or not a given nexus has been announced.
func (l *BaseLayer) IsNexusAnnounced(nexus nexusHelper.Nexus) bool {
	if _, ok := l.Announced.Get(nexus.UUID()); ok {
		return true
	}
	return false
}

// SendLayerEvent sends an event up the chain.
func (l *BaseLayer) SendLayerEvent(nexus nexusHelper.Nexus, status string) {
	if l.OnLayerEvent != nil {
		l.OnLayerEvent(l.LayerName, nexus, status)
	}
}

// TrackNexusRevoked tracks that a nexus has been revoked.
func (l *BaseLayer) TrackNexusRevoked(nexus nexusHelper.Nexus) {
	if l.IsNexusAnnounced(nexus) {
		l.Announced.Delete(nexus.UUID())
	}
}

// RevokeNexus removes the nexus from directories/layers.
func (l *BaseLayer) RevokeNexus(belowNexus nexusHelper.Nexus) {
	belowUUID, _ := l.NexusByBelow.Get(belowNexus.UUID())
	nexus, _ := l.Nexuses.Get(belowUUID)
	l.UntrackNexus(nexus, belowNexus)
	if l.IsNexusAnnounced(nexus) {
		l.TrackNexusRevoked(nexus)
		if l.OnRevoke != nil {
			l.OnRevoke(nexus)
		}
		l.SendLayerEvent(nexus, message.NexusRevoked)
	}
}

// NewBaseLayer creates a new base layer, note you still have to call register_above nexus.
func NewBaseLayer() BaseLayer {
	return BaseLayer{
		Nexuses:      NexusMap{},
		Announced:    NexusMap{},
		BelowNexuses: NexusMap{},
		NexusByBelow: NexusUUIDMap{},
		BelowByNexus: NexusUUIDMap{},
	}
}
