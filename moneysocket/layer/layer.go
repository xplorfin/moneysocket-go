package layer

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

type OnLayerEventFn = func(layerName string, nexus nexusHelper.Nexus, event string)
type OnAnnounceFn = func(nexus nexusHelper.Nexus)
type OnRevokeFn = func(nexus nexusHelper.Nexus)

type Layer interface {
	// set on layer event
	SetOnLayerEvent(o OnLayerEventFn)
	// set on announce event
	SetOnAnnounce(o OnAnnounceFn)
	// set on revoke event
	SetOnRevoke(o OnAnnounceFn)
	// register above layer events with current layer
	// must be done here since announce nexus
	// is not available form base layer
	RegisterAboveLayer(belowLayer Layer)
	// register layer event with nexuses
	RegisterLayerEvent(fn OnLayerEventFn, layerName string)
	// announce nexusHelper
	AnnounceNexus(belowNexus nexusHelper.Nexus)
	// track nexusHelper events
	TrackNexus(nexus nexusHelper.Nexus, belowNexus nexusHelper.Nexus)
	// remove tracker from nexusHelper events
	UntrackNexus(nexus nexusHelper.Nexus, belowNexus nexusHelper.Nexus)
	// track that a nexus has been announced
	TrackNexusAnnounced(nexus nexusHelper.Nexus)
	// check wether or not a nexus has been announced
	IsNexusAnnounced(nexus nexusHelper.Nexus) bool
	// revoke a nexus
	RevokeNexus(belowNexus nexusHelper.Nexus)
}

type BaseLayer struct {
	LayerName    string
	OnLayerEvent OnLayerEventFn
	OnAnnounce   OnAnnounceFn
	OnRevoke     OnRevokeFn

	Nexuses      NexusMap
	Announced    NexusMap
	BelowNexuses NexusMap
	NexusByBelow NexusUuidMap
	BelowByNexus NexusUuidMap
}

func (l *BaseLayer) SetOnLayerEvent(o OnLayerEventFn) {
	l.OnLayerEvent = o
}

func (l *BaseLayer) SetOnAnnounce(o OnAnnounceFn) {
	l.OnAnnounce = o
}

func (l *BaseLayer) SetOnRevoke(o OnRevokeFn) {
	l.OnRevoke = o
}

func (l *BaseLayer) RegisterLayerEvent(fn OnLayerEventFn, layerName string) {
	l.LayerName = layerName
	l.OnLayerEvent = fn
}

func (l *BaseLayer) TrackNexus(nexus nexusHelper.Nexus, belowNexus nexusHelper.Nexus) {
	l.Nexuses.Store(nexus.Uuid(), nexus)
	l.BelowNexuses.Store(belowNexus.Uuid(), belowNexus)
	l.NexusByBelow.Store(belowNexus.Uuid(), nexus.Uuid())
	l.BelowByNexus.Store(nexus.Uuid(), belowNexus.Uuid())
	l.SendLayerEvent(nexus, message.NexusCreated)
}

func (l *BaseLayer) UntrackNexus(nexus nexusHelper.Nexus, belowNexus nexusHelper.Nexus) {
	l.Nexuses.Delete(nexus.Uuid())
	l.BelowByNexus.Delete(belowNexus.Uuid())
	l.NexusByBelow.Delete(belowNexus.Uuid())
	l.BelowByNexus.Delete(belowNexus.Uuid())
	l.SendLayerEvent(nexus, message.NexusDestroyed)
}

func (l *BaseLayer) TrackNexusAnnounced(nexus nexusHelper.Nexus) {
	l.Announced.Store(nexus.Uuid(), nexus)
}

func (l *BaseLayer) IsNexusAnnounced(nexus nexusHelper.Nexus) bool {
	if _, ok := l.Announced.Get(nexus.Uuid()); ok {
		return true
	}
	return false
}

func (l *BaseLayer) SendLayerEvent(nexus nexusHelper.Nexus, status string) {
	if l.OnLayerEvent != nil {
		l.OnLayerEvent(l.LayerName, nexus, status)
	}
}

func (l *BaseLayer) TrackNexusRevoked(nexus nexusHelper.Nexus) {
	if l.IsNexusAnnounced(nexus) {
		l.Announced.Delete(nexus.Uuid())
	}
}

func (l *BaseLayer) RevokeNexus(belowNexus nexusHelper.Nexus) {
	belowUuid, _ := l.NexusByBelow.Get(belowNexus.Uuid())
	nexus, _ := l.Nexuses.Get(belowUuid)
	l.UntrackNexus(nexus, belowNexus)
	if l.IsNexusAnnounced(nexus) {
		l.TrackNexusRevoked(nexus)
		if l.OnRevoke != nil {
			l.OnRevoke(nexus)
		}
		l.SendLayerEvent(nexus, message.NexusRevoked)
	}
}

// create a new base layer, note you still have to call register_above nexus
func NewBaseLayer() BaseLayer {
	return BaseLayer{
		Nexuses:      NexusMap{},
		Announced:    NexusMap{},
		BelowNexuses: NexusMap{},
		NexusByBelow: NexusUuidMap{},
		BelowByNexus: NexusUuidMap{},
	}
}
