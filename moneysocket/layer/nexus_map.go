package layer

import (
	"sync"

	uuid "github.com/satori/go.uuid"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

// NexusMap that enforces types
// TODO this should get replaced with a string map
type NexusMap struct {
	// map[uuid.UUID]nexusHelper.Nexus
	internalMap NexusStringMap
}

// Store idempotently add an nexus to map
func (nm *NexusMap) Store(nexusUUID uuid.UUID, nexus nexusHelper.Nexus) {
	nm.internalMap.Store(nexusUUID.String(), nexus)
}

// Delete deletes an item from a map
func (nm *NexusMap) Delete(nexusUUID uuid.UUID) {
	nm.internalMap.Delete(nexusUUID.String())
}

// Get a nexus from the map, don't check for key presence
func (nm *NexusMap) Get(nexusUUID uuid.UUID) (nexusHelper.Nexus, bool) {
	nu, ok := nm.internalMap.Get(nexusUUID.String())
	if !ok {
		return NewUnknownNexus(), ok
	}
	return nu.(nexusHelper.Nexus), ok
}

// Range loops through a nexus map
func (nm *NexusMap) Range(f func(key uuid.UUID, nexus nexusHelper.Nexus) bool) {
	nm.internalMap.Range(func(key string, value nexusHelper.Nexus) bool {
		return f(uuid.FromStringOrNil(key), value.(nexusHelper.Nexus))
	})
}

// NexusUUIDMap is the nexus map that enforces types
type NexusUUIDMap struct {
	// map[uuid.UUID]uuid.UUID
	internalMap sync.Map
}

// Store idempotently add an idea to map
func (num *NexusUUIDMap) Store(nexusUUID uuid.UUID, nexus uuid.UUID) {
	num.internalMap.Store(nexusUUID, nexus)
}

// Delete an item from a map
func (num *NexusUUIDMap) Delete(nexusUUID uuid.UUID) {
	num.internalMap.Delete(nexusUUID)
}

// Get a nexus from the map, don't check for key presence
func (num *NexusUUIDMap) Get(nexusUUID uuid.UUID) (uuid.UUID, bool) {
	nu, ok := num.internalMap.Load(nexusUUID)
	if !ok {
		return uuid.NewV4(), ok
	}
	return nu.(uuid.UUID), ok
}

// Range iterates through a NexusUUIDMap
func (num *NexusUUIDMap) Range(f func(key uuid.UUID, value uuid.UUID) bool) {
	num.internalMap.Range(func(key, value interface{}) bool {
		return f(key.(uuid.UUID), value.(uuid.UUID))
	})
}

// NexusStringMap is a nexus map that enforces types
type NexusStringMap struct {
	// map[uuid.UUID]nexus
	internalMap sync.Map
}

// Store idempotently adds an item to map
func (num *NexusStringMap) Store(key string, nexus nexusHelper.Nexus) {
	num.internalMap.Store(key, nexus)
}

// Delete an item from a map
func (num *NexusStringMap) Delete(key string) {
	num.internalMap.Delete(key)
}

// Get a nexus from the map, don't check for key presence
func (num *NexusStringMap) Get(key string) (nexusHelper.Nexus, bool) {
	nu, ok := num.internalMap.Load(key)
	if !ok {
		return NewUnknownNexus(), ok
	}
	return nu.(nexusHelper.Nexus), ok
}

// Range loops through a NexusStringMap
func (num *NexusStringMap) Range(f func(key string, value nexusHelper.Nexus) bool) {
	num.internalMap.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(nexusHelper.Nexus))
	})
}
