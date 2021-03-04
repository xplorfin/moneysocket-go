package rendezvous

import (
	"sync"

	uuid "github.com/satori/go.uuid"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

type RendezvousDirectory struct {
	// todo switch to sync maps with individual mutexes
	mux sync.Mutex
	// map from uuid to cooresponding nexus
	NexusesByUUID map[uuid.UUID]nexusHelper.Nexus
	// map from nexus uuid to cooresponding rendezvous id
	RidsByUUIDs map[uuid.UUID]string
	// map of uuids that are currently unpaired sorted by rendezvous id
	UnpeeredUUIDsByRid map[string]uuid.UUID
	// map of nexus to the uuid of the cooresponding peer
	UuidPeers map[uuid.UUID]uuid.UUID
	// list of rids which are currently peered
	RidsPeered []string
}

func NewRendezvousDirectory() *RendezvousDirectory {
	return &RendezvousDirectory{
		NexusesByUUID:      make(map[uuid.UUID]nexusHelper.Nexus),
		RidsByUUIDs:        make(map[uuid.UUID]string),
		UnpeeredUUIDsByRid: make(map[string]uuid.UUID),
		UuidPeers:          make(map[uuid.UUID]uuid.UUID),
		RidsPeered:         []string{},
	}
}

// get matching peer of a given nexus (1:1 mappnig)
func (r *RendezvousDirectory) GetPeerNexus(nexusId uuid.UUID) *nexusHelper.Nexus {
	if peer, ok := r.UuidPeers[nexusId]; ok {
		matchingN := r.NexusesByUUID[peer]
		return &matchingN
	}
	return nil
}

// wether or not nexus is indexed in directory
func (r *RendezvousDirectory) IsNexusInDirectory(nexus nexusHelper.Nexus) bool {
	if _, ok := r.NexusesByUUID[nexus.Uuid()]; ok {
		return true
	}
	return false
}

// check if a rendezvous id currently has a peer
func (r *RendezvousDirectory) IsRidPeered(rendezvousId string) bool {
	for _, rid := range r.RidsPeered {
		if rid == rendezvousId {
			return true
		}
	}
	return false
}

// add and index a nexus
func (r *RendezvousDirectory) AddNexus(nexus nexusHelper.Nexus, rendezvousId string) {
	r.mux.Lock()
	r.NexusesByUUID[nexus.Uuid()] = nexus
	r.RidsByUUIDs[nexus.Uuid()] = rendezvousId
	if peerUuid, ok := r.UnpeeredUUIDsByRid[rendezvousId]; ok {
		delete(r.UnpeeredUUIDsByRid, rendezvousId)
		r.UuidPeers[peerUuid] = nexus.Uuid()
		r.UuidPeers[nexus.Uuid()] = peerUuid
		r.RidsPeered = append(r.RidsPeered, rendezvousId)
	} else {
		r.UnpeeredUUIDsByRid[rendezvousId] = nexus.Uuid()
	}
	r.mux.Unlock()
}

// remove a nexus from the directory
func (r *RendezvousDirectory) RemoveNexus(nexus nexusHelper.Nexus) {
	if !r.IsNexusInDirectory(nexus) {
		return
	}
	r.mux.Lock()
	delete(r.NexusesByUUID, nexus.Uuid())
	rid := r.RidsByUUIDs[nexus.Uuid()]
	delete(r.RidsByUUIDs, nexus.Uuid())

	if peerUuid, ok := r.UuidPeers[nexus.Uuid()]; ok {
		delete(r.UuidPeers, peerUuid)
		delete(r.UuidPeers, nexus.Uuid())
		r.UnpeeredUUIDsByRid[rid] = peerUuid
		// remove from peered uuids
		for i, peeredRid := range r.RidsPeered {
			if peeredRid == rid {
				r.RidsPeered = append(r.RidsPeered[:i], r.RidsPeered[i+1:]...)
			}
		}
	} else {
		delete(r.UnpeeredUUIDsByRid, rid)
	}
	r.mux.Unlock()
}
