package rendezvous

import (
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

// Directory stores nexus peering data
type Directory struct {
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
func NewRendezvousDirectory() *Directory {
	return &Directory{
		NexusesByUUID:      make(map[string]nexusHelper.Nexus),
		RidsByUUIDs:        make(map[string]string),
		UnpeeredUUIDsByRid: make(map[string]string),
		UUIDPeers:          make(map[string]string),
		RIDSPeered:         []string{},
	}
}

// ToString gets the number of nexuses/peered/unpeered
func (r *Directory) ToString() string {
	nexuses := len(r.NexusesByUUID)
	unpeered := len(r.UnpeeredUUIDsByRid)
	peered := len(r.UUIDPeers)
	return fmt.Sprintf("nexuses/unpeered/peered %d/%d/%d", nexuses, unpeered, peered)
}

// GetPeerNexus get matching peer of a given nexus (1:1 mapping)
func (r *Directory) GetPeerNexus(nexusID uuid.UUID) *nexusHelper.Nexus {
	if peer, ok := r.UUIDPeers[nexusID.String()]; ok {
		matchingN := r.NexusesByUUID[peer]
		return &matchingN
	}
	return nil
}

// IsNexusInDirectory determines wether or not nexus is indexed in directory
func (r *Directory) IsNexusInDirectory(nexus nexusHelper.Nexus) bool {
	if _, ok := r.NexusesByUUID[nexus.UUID().String()]; ok {
		return true
	}
	return false
}

// IsRidPeered checks if a rendezvous id currently has a peer
func (r *Directory) IsRidPeered(rendezvousID string) bool {
	for _, rid := range r.RIDSPeered {
		if rid == rendezvousID {
			return true
		}
	}
	return false
}

// AddNexus adds and indexes a nexus
func (r *Directory) AddNexus(nexus nexusHelper.Nexus, rendezvousID string) {
	r.mux.Lock()
	fmt.Print(nexus.UUID().String())
	r.NexusesByUUID[nexus.UUID().String()] = nexus
	r.RidsByUUIDs[nexus.UUID().String()] = rendezvousID
	if peerUUID, ok := r.UnpeeredUUIDsByRid[rendezvousID]; ok {
		delete(r.UnpeeredUUIDsByRid, rendezvousID)
		r.UUIDPeers[peerUUID] = nexus.UUID().String()
		r.UUIDPeers[nexus.UUID().String()] = peerUUID
		r.RIDSPeered = append(r.RIDSPeered, rendezvousID)
	} else {
		r.UnpeeredUUIDsByRid[rendezvousID] = nexus.UUID().String()
	}
	r.mux.Unlock()
}

// RemoveNexus removes a nexus from the directory
func (r *Directory) RemoveNexus(nexus nexusHelper.Nexus) {
	if !r.IsNexusInDirectory(nexus) {
		return
	}
	r.mux.Lock()
	delete(r.NexusesByUUID, nexus.UUID().String())
	rid := r.RidsByUUIDs[nexus.UUID().String()]
	delete(r.RidsByUUIDs, nexus.UUID().String())

	if peerUUID, ok := r.UUIDPeers[nexus.UUID().String()]; ok {
		delete(r.UUIDPeers, peerUUID)
		delete(r.UUIDPeers, nexus.UUID().String())
		r.UnpeeredUUIDsByRid[rid] = peerUUID
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
