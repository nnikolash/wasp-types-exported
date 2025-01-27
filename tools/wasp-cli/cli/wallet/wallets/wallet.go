package wallets

import "github.com/nnikolash/wasp-types-exported/packages/cryptolib"

type Wallet interface {
	cryptolib.VariantKeyPair

	AddressIndex() uint32
}
