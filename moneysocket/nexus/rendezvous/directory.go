package rendezvous

import (
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

// RendezvousDirectory stores nexus peering data
type RendezvousDirectory struct {
	// todo switch to sync maps with individual mutexes
	// this can probably be removed rn
	mux sync.Mutex
	// NexusesByUUID is a map from uuid to corresponding nexus
	NexusesByUUID map[string]nexusHelper.Nexus
	// RidsByUUIDs is a map from nexus uuid to corresponding rendezvous id
	RidsByUUIDs map[string]string
	// UnpeeredUUIDsByRid is a map of uuids that are currently unpaired sorted by rendezvous id
	UnpeeredUUIDsByRid map[string]string
	// UUIDPeers is a map of nexus to the uuid of the corresponding peer
	UUIDPeers map[string]string
	// RIDSPeered is a list of rids which are currently peered
	RIDSPeered []string
}

// NewRendezvousDirectory creates a new RendezvousDirectory
func NewRendezvousDirectory() *RendezvousDirectory {
	return &RendezvousDirectory{
		NexusesByUUID:      make(map[string]nexusHelper.Nexus),
		RidsByUUIDs:        make(map[string]string),
		UnpeeredUUIDsByRid: make(map[string]string),
		UUIDPeers:          make(map[string]string),
		RIDSPeered:         []string{},
	}
}

// ToString gets the number of nexuses/peered/unpeered
func (r *RendezvousDirectory) ToString() string {
	nexuses := len(r.NexusesByUUID)
	unpeered := len(r.UnpeeredUUIDsByRid)
	peered := len(r.UUIDPeers)
	return fmt.Sprintf("nexuses/unpeered/peered %d/%d/%d", nexuses, unpeered, peered)
}

// GetPeerNexus get matching peer of a given nexus (1:1 mapping)
func (r *RendezvousDirectory) GetPeerNexus(nexusId uuid.UUID) *nexusHelper.Nexus {
	if peer, ok := r.UUIDPeers[nexusId.String()]; ok {
		matchingN := r.NexusesByUUID[peer]
		return &matchingN
	}
	return nil
}

// IsNexusInDirectory determines wether or not nexus is indexed in directory
func (r *RendezvousDirectory) IsNexusInDirectory(nexus nexusHelper.Nexus) bool {
	if _, ok := r.NexusesByUUID[nexus.UUID().String()]; ok {
		return true
	}
	return false
}

// IsRidPeered checks if a rendezvous id currently has a peer
func (r *RendezvousDirectory) IsRidPeered(rendezvousId string) bool {
	for _, rid := range r.RIDSPeered {
		if rid == rendezvousId {
			return true
		}
	}
	return false
}

// AddNexus adds and indexes a nexus
func (r *RendezvousDirectory) AddNexus(nexus nexusHelper.Nexus, rendezvousId string) {
	r.mux.Lock()
	fmt.Print(nexus.UUID().String())
	r.NexusesByUUID[nexus.UUID().String()] = nexus
	r.RidsByUUIDs[nexus.UUID().String()] = rendezvousId
	if peerUuid, ok := r.UnpeeredUUIDsByRid[rendezvousId]; ok {
		delete(r.UnpeeredUUIDsByRid, rendezvousId)
		r.UUIDPeers[peerUuid] = nexus.UUID().String()
		r.UUIDPeers[nexus.UUID().String()] = peerUuid
		r.RIDSPeered = append(r.RIDSPeered, rendezvousId)
	} else {
		r.UnpeeredUUIDsByRid[rendezvousId] = nexus.UUID().String()
	}
	r.mux.Unlock()
}

// RemoveNexus removes a nexus from the directory
func (r *RendezvousDirectory) RemoveNexus(nexus nexusHelper.Nexus) {
	if !r.IsNexusInDirectory(nexus) {
		return
	}
	r.mux.Lock()
	delete(r.NexusesByUUID, nexus.UUID().String())
	rid := r.RidsByUUIDs[nexus.UUID().String()]
	delete(r.RidsByUUIDs, nexus.UUID().String())

	if peerUuid, ok := r.UUIDPeers[nexus.UUID().String()]; ok {
		delete(r.UUIDPeers, peerUuid)
		delete(r.UUIDPeers, nexus.UUID().String())
		r.UnpeeredUUIDsByRid[rid] = peerUuid
		// remove from peered uuids
		for i, peeredRid := range r.RIDSPeered {
			if peeredRid == rid {
				r.RIDSPeered = append(r.RIDSPeered[:i], r.RIDSPeered[i+1:]...)
			}
		}
	} else {
		delete(r.UnpeeredUUIDsByRid, rid)
	}
	r.mux.Unlock()
}
