package isc

import (
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"

	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

type ContractIdentityKind rwutil.Kind

type ContractIdentity struct {
	// can either be an Hname or a solidity contract
	Kind ContractIdentityKind

	// only 1 or the other will be filled
	EvmAddr common.Address
	hname   Hname
}

const (
	ContractIdentityKindEmpty ContractIdentityKind = iota
	ContractIdentityKindHname
	ContractIdentityKindEthereum
)

func EmptyContractIdentity() ContractIdentity {
	return ContractIdentity{Kind: ContractIdentityKindEmpty}
}

func ContractIdentityFromHname(hn Hname) ContractIdentity {
	return ContractIdentity{hname: hn, Kind: ContractIdentityKindHname}
}

func ContractIdentityFromEVMAddress(addr common.Address) ContractIdentity {
	return ContractIdentity{EvmAddr: addr, Kind: ContractIdentityKindEthereum}
}

func (c *ContractIdentity) String() string {
	switch c.Kind {
	case ContractIdentityKindHname:
		return c.hname.String()
	case ContractIdentityKindEthereum:
		return c.EvmAddr.String()
	}
	return ""
}

func (c *ContractIdentity) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	c.Kind = ContractIdentityKind(rr.ReadKind())
	switch c.Kind {
	case ContractIdentityKindHname:
		rr.Read(&c.hname)
	case ContractIdentityKindEthereum:
		rr.ReadN(c.EvmAddr[:])
	}
	return rr.Err
}

func (c *ContractIdentity) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteKind(rwutil.Kind(c.Kind))
	switch c.Kind {
	case ContractIdentityKindHname:
		ww.Write(&c.hname)
	case ContractIdentityKindEthereum:
		ww.WriteN(c.EvmAddr[:])
	}
	return ww.Err
}

func (c *ContractIdentity) AgentID(chainID ChainID) AgentID {
	switch c.Kind {
	case ContractIdentityKindHname:
		return NewContractAgentID(chainID, c.hname)
	case ContractIdentityKindEthereum:
		return NewEthereumAddressAgentID(chainID, c.EvmAddr)
	}
	return &NilAgentID{}
}

func (c *ContractIdentity) Hname() (Hname, error) {
	if c.Kind == ContractIdentityKindHname {
		return c.hname, nil
	}
	return 0, fmt.Errorf("not an Hname contract")
}

func (c *ContractIdentity) HnameRaw() Hname {
	return c.hname
}

func (c *ContractIdentity) EVMAddress() (common.Address, error) {
	if c.Kind == ContractIdentityKindEthereum {
		return c.EvmAddr, nil
	}
	return common.Address{}, fmt.Errorf("not an EVM contract")
}

func (c *ContractIdentity) Empty() bool {
	return c.Kind == ContractIdentityKindEmpty
}
