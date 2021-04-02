package compat

import "github.com/xplorfin/moneysocket-go/moneysocket/nexus"

// RevokableNexus is a nexus that can be revoked from a layer.
type RevokableNexus interface {
	nexus.Nexus
	// RevokeFromLayer revokes a nexus from a given layer
	RevokeFromLayer()
}
