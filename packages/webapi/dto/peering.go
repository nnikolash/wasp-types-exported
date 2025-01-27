package dto

import "github.com/nnikolash/wasp-types-exported/packages/cryptolib"

type PeeringNodeStatus struct {
	Name       string
	IsAlive    bool
	PeeringURL string
	NumUsers   int
	PublicKey  *cryptolib.PublicKey
	IsTrusted  bool
}

type PeeringNodeIdentity struct {
	Name       string
	PublicKey  *cryptolib.PublicKey
	PeeringURL string
	IsTrusted  bool
}
