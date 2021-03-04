package layer

import (
	"sync"

	uuid "github.com/satori/go.uuid"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

// nexus map that enforces
// TODO this should get replaced with a string map
type NexusMap struct {
	// map[uuid.UUID]nexusHelper.Nexus
	internalMap NexusStringMap
}

// idempotently add an idea to map
func (nm *NexusMap) Store(nexusUuid uuid.UUID, nexus nexusHelper.Nexus) {
	nm.internalMap.Store(nexusUuid.String(), nexus)
}

// delete an item from a map
func (nm *NexusMap) Delete(nexusUuid uuid.UUID) {
	nm.internalMap.Delete(nexusUuid.String())
}

// get a nexus from the map, don't check for key presence
func (nm *NexusMap) Get(nexusUuid uuid.UUID) (nexusHelper.Nexus, bool) {
	nu, ok := nm.internalMap.Get(nexusUuid.String())
	if !ok {
		return NewUnknownNexus(), ok
	}
	return nu.(nexusHelper.Nexus), ok
}

func (nm *NexusMap) Range(f func(key uuid.UUID, nexus nexusHelper.Nexus) bool) {
	nm.internalMap.Range(func(key string, value nexusHelper.Nexus) bool {
		return f(uuid.FromStringOrNil(key), value.(nexusHelper.Nexus))
	})
}

// nexus map that enforces types
type NexusUuidMap struct {
	// map[uuid.UUID]uuid.UUID
	internalMap sync.Map
}

// idempotently add an idea to map
func (num *NexusUuidMap) Store(nexusUuid uuid.UUID, nexus uuid.UUID) {
	num.internalMap.Store(nexusUuid, nexus)
}

// delete an item from a map
func (num *NexusUuidMap) Delete(nexusUuid uuid.UUID) {
	num.internalMap.Delete(nexusUuid)
}

// get a nexus from the map, don't check for key presence
func (num *NexusUuidMap) Get(nexusUuid uuid.UUID) (uuid.UUID, bool) {
	nu, ok := num.internalMap.Load(nexusUuid)
	if !ok {
		return uuid.NewV4(), ok
	}
	return nu.(uuid.UUID), ok
}

func (num *NexusUuidMap) Range(f func(key uuid.UUID, value uuid.UUID) bool) {
	num.internalMap.Range(func(key, value interface{}) bool {
		return f(key.(uuid.UUID), value.(uuid.UUID))
	})
}

// nexus map that enforces types
type NexusStringMap struct {
	// map[uuid.UUID]nexus
	internalMap sync.Map
}

// idempotently add an idea to map
func (num *NexusStringMap) Store(key string, nexus nexusHelper.Nexus) {
	num.internalMap.Store(key, nexus)
}

// delete an item from a map
func (num *NexusStringMap) Delete(key string) {
	num.internalMap.Delete(key)
}

// get a nexus from the map, don't check for key presence
func (num *NexusStringMap) Get(key string) (nexusHelper.Nexus, bool) {
	nu, ok := num.internalMap.Load(key)
	if !ok {
		return NewUnknownNexus(), ok
	}
	return nu.(nexusHelper.Nexus), ok
}

func (num *NexusStringMap) Range(f func(key string, value nexusHelper.Nexus) bool) {
	num.internalMap.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(nexusHelper.Nexus))
	})
}
