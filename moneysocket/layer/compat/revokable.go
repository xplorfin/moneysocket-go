package compat

import "github.com/xplorfin/moneysocket-go/moneysocket/nexus"

// revoke from a layer
type RevokableNexus interface {
	nexus.Nexus
	RevokeFromLayer()
}
