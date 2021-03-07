package rendezvous

import (
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

type RendezvousDirectory struct {
	// todo switch to sync maps with individual mutexes
	mux sync.Mutex
	// map from uuid to cooresponding nexus
	NexusesByUUID map[string]nexusHelper.Nexus
	// map from nexus uuid to cooresponding rendezvous id
	RidsByUUIDs map[string]string
	// map of uuids that are currently unpaired sorted by rendezvous id
	UnpeeredUUIDsByRid map[string]string
	// map of nexus to the uuid of the cooresponding peer
	UuidPeers map[string]string
	// list of rids which are currently peered
	RidsPeered []string
}

func NewRendezvousDirectory() *RendezvousDirectory {
	return &RendezvousDirectory{
		NexusesByUUID:      make(map[string]nexusHelper.Nexus),
		RidsByUUIDs:        make(map[string]string),
		UnpeeredUUIDsByRid: make(map[string]string),
		UuidPeers:          make(map[string]string),
		RidsPeered:         []string{},
	}
}

// get matching peer of a given nexus (1:1 mappnig)
func (r *RendezvousDirectory) GetPeerNexus(nexusId uuid.UUID) *nexusHelper.Nexus {
	if peer, ok := r.UuidPeers[nexusId.String()]; ok {
		matchingN := r.NexusesByUUID[peer]
		return &matchingN
	}
	return nil
}

// wether or not nexus is indexed in directory
func (r *RendezvousDirectory) IsNexusInDirectory(nexus nexusHelper.Nexus) bool {
	if _, ok := r.NexusesByUUID[nexus.Uuid().String()]; ok {
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
	fmt.Print(nexus.Uuid().String())
	r.NexusesByUUID[nexus.Uuid().String()] = nexus
	r.RidsByUUIDs[nexus.Uuid().String()] = rendezvousId
	if peerUuid, ok := r.UnpeeredUUIDsByRid[rendezvousId]; ok {
		delete(r.UnpeeredUUIDsByRid, rendezvousId)
		r.UuidPeers[peerUuid] = nexus.Uuid().String()
		r.UuidPeers[nexus.Uuid().String()] = peerUuid
		r.RidsPeered = append(r.RidsPeered, rendezvousId)
	} else {
		r.UnpeeredUUIDsByRid[rendezvousId] = nexus.Uuid().String()
	}
	r.mux.Unlock()
}

// remove a nexus from the directory
func (r *RendezvousDirectory) RemoveNexus(nexus nexusHelper.Nexus) {
	if !r.IsNexusInDirectory(nexus) {
		return
	}
	r.mux.Lock()
	delete(r.NexusesByUUID, nexus.Uuid().String())
	rid := r.RidsByUUIDs[nexus.Uuid().String()]
	delete(r.RidsByUUIDs, nexus.Uuid().String())

	if peerUuid, ok := r.UuidPeers[nexus.Uuid().String()]; ok {
		delete(r.UuidPeers, peerUuid)
		delete(r.UuidPeers, nexus.Uuid().String())
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
