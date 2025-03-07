// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

#![allow(dead_code)]
#![allow(unused_imports)]

use crate::*;
use crate::coreaccounts::*;

#[derive(Clone)]
pub struct ImmutableFoundryCreateNewParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableFoundryCreateNewParams {
    pub fn new() -> ImmutableFoundryCreateNewParams {
        ImmutableFoundryCreateNewParams {
            proxy: params_proxy(),
        }
    }

    // token scheme for the new foundry
    pub fn token_scheme(&self) -> ScImmutableBytes {
        ScImmutableBytes::new(self.proxy.root(PARAM_TOKEN_SCHEME))
    }
}

#[derive(Clone)]
pub struct MutableFoundryCreateNewParams {
    pub(crate) proxy: Proxy,
}

impl MutableFoundryCreateNewParams {
    // token scheme for the new foundry
    pub fn token_scheme(&self) -> ScMutableBytes {
        ScMutableBytes::new(self.proxy.root(PARAM_TOKEN_SCHEME))
    }
}

#[derive(Clone)]
pub struct ImmutableNativeTokenCreateParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableNativeTokenCreateParams {
    pub fn new() -> ImmutableNativeTokenCreateParams {
        ImmutableNativeTokenCreateParams {
            proxy: params_proxy(),
        }
    }

    pub fn token_decimals(&self) -> ScImmutableUint8 {
        ScImmutableUint8::new(self.proxy.root(PARAM_TOKEN_DECIMALS))
    }

    pub fn token_name(&self) -> ScImmutableString {
        ScImmutableString::new(self.proxy.root(PARAM_TOKEN_NAME))
    }

    // token scheme for the new foundry
    pub fn token_scheme(&self) -> ScImmutableBytes {
        ScImmutableBytes::new(self.proxy.root(PARAM_TOKEN_SCHEME))
    }

    pub fn token_symbol(&self) -> ScImmutableString {
        ScImmutableString::new(self.proxy.root(PARAM_TOKEN_SYMBOL))
    }
}

#[derive(Clone)]
pub struct MutableNativeTokenCreateParams {
    pub(crate) proxy: Proxy,
}

impl MutableNativeTokenCreateParams {
    pub fn token_decimals(&self) -> ScMutableUint8 {
        ScMutableUint8::new(self.proxy.root(PARAM_TOKEN_DECIMALS))
    }

    pub fn token_name(&self) -> ScMutableString {
        ScMutableString::new(self.proxy.root(PARAM_TOKEN_NAME))
    }

    // token scheme for the new foundry
    pub fn token_scheme(&self) -> ScMutableBytes {
        ScMutableBytes::new(self.proxy.root(PARAM_TOKEN_SCHEME))
    }

    pub fn token_symbol(&self) -> ScMutableString {
        ScMutableString::new(self.proxy.root(PARAM_TOKEN_SYMBOL))
    }
}

#[derive(Clone)]
pub struct ImmutableNativeTokenDestroyParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableNativeTokenDestroyParams {
    pub fn new() -> ImmutableNativeTokenDestroyParams {
        ImmutableNativeTokenDestroyParams {
            proxy: params_proxy(),
        }
    }

    // serial number of the foundry
    pub fn foundry_sn(&self) -> ScImmutableUint32 {
        ScImmutableUint32::new(self.proxy.root(PARAM_FOUNDRY_SN))
    }
}

#[derive(Clone)]
pub struct MutableNativeTokenDestroyParams {
    pub(crate) proxy: Proxy,
}

impl MutableNativeTokenDestroyParams {
    // serial number of the foundry
    pub fn foundry_sn(&self) -> ScMutableUint32 {
        ScMutableUint32::new(self.proxy.root(PARAM_FOUNDRY_SN))
    }
}

#[derive(Clone)]
pub struct ImmutableNativeTokenModifySupplyParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableNativeTokenModifySupplyParams {
    pub fn new() -> ImmutableNativeTokenModifySupplyParams {
        ImmutableNativeTokenModifySupplyParams {
            proxy: params_proxy(),
        }
    }

    // mint (default) or destroy tokens
    pub fn destroy_tokens(&self) -> ScImmutableBool {
        ScImmutableBool::new(self.proxy.root(PARAM_DESTROY_TOKENS))
    }

    // serial number of the foundry
    pub fn foundry_sn(&self) -> ScImmutableUint32 {
        ScImmutableUint32::new(self.proxy.root(PARAM_FOUNDRY_SN))
    }

    // positive nonzero amount to mint or destroy
    pub fn supply_delta_abs(&self) -> ScImmutableBigInt {
        ScImmutableBigInt::new(self.proxy.root(PARAM_SUPPLY_DELTA_ABS))
    }
}

#[derive(Clone)]
pub struct MutableNativeTokenModifySupplyParams {
    pub(crate) proxy: Proxy,
}

impl MutableNativeTokenModifySupplyParams {
    // mint (default) or destroy tokens
    pub fn destroy_tokens(&self) -> ScMutableBool {
        ScMutableBool::new(self.proxy.root(PARAM_DESTROY_TOKENS))
    }

    // serial number of the foundry
    pub fn foundry_sn(&self) -> ScMutableUint32 {
        ScMutableUint32::new(self.proxy.root(PARAM_FOUNDRY_SN))
    }

    // positive nonzero amount to mint or destroy
    pub fn supply_delta_abs(&self) -> ScMutableBigInt {
        ScMutableBigInt::new(self.proxy.root(PARAM_SUPPLY_DELTA_ABS))
    }
}

#[derive(Clone)]
pub struct ImmutableTransferAccountToChainParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableTransferAccountToChainParams {
    pub fn new() -> ImmutableTransferAccountToChainParams {
        ImmutableTransferAccountToChainParams {
            proxy: params_proxy(),
        }
    }

    // Optional gas amount to reserve in the allowance for the internal
    // call to transferAllowanceTo(). Default 10_000 (MinGasFee).
    pub fn gas_reserve(&self) -> ScImmutableUint64 {
        ScImmutableUint64::new(self.proxy.root(PARAM_GAS_RESERVE))
    }
}

#[derive(Clone)]
pub struct MutableTransferAccountToChainParams {
    pub(crate) proxy: Proxy,
}

impl MutableTransferAccountToChainParams {
    // Optional gas amount to reserve in the allowance for the internal
    // call to transferAllowanceTo(). Default 10_000 (MinGasFee).
    pub fn gas_reserve(&self) -> ScMutableUint64 {
        ScMutableUint64::new(self.proxy.root(PARAM_GAS_RESERVE))
    }
}

#[derive(Clone)]
pub struct ImmutableTransferAllowanceToParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableTransferAllowanceToParams {
    pub fn new() -> ImmutableTransferAllowanceToParams {
        ImmutableTransferAllowanceToParams {
            proxy: params_proxy(),
        }
    }

    // The target L2 account
    pub fn agent_id(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct MutableTransferAllowanceToParams {
    pub(crate) proxy: Proxy,
}

impl MutableTransferAllowanceToParams {
    // The target L2 account
    pub fn agent_id(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct ImmutableAccountFoundriesParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableAccountFoundriesParams {
    pub fn new() -> ImmutableAccountFoundriesParams {
        ImmutableAccountFoundriesParams {
            proxy: params_proxy(),
        }
    }

    // account agent ID
    pub fn agent_id(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct MutableAccountFoundriesParams {
    pub(crate) proxy: Proxy,
}

impl MutableAccountFoundriesParams {
    // account agent ID
    pub fn agent_id(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct ImmutableAccountNFTAmountParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableAccountNFTAmountParams {
    pub fn new() -> ImmutableAccountNFTAmountParams {
        ImmutableAccountNFTAmountParams {
            proxy: params_proxy(),
        }
    }

    // account agent ID
    pub fn agent_id(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct MutableAccountNFTAmountParams {
    pub(crate) proxy: Proxy,
}

impl MutableAccountNFTAmountParams {
    // account agent ID
    pub fn agent_id(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct ImmutableAccountNFTAmountInCollectionParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableAccountNFTAmountInCollectionParams {
    pub fn new() -> ImmutableAccountNFTAmountInCollectionParams {
        ImmutableAccountNFTAmountInCollectionParams {
            proxy: params_proxy(),
        }
    }

    // account agent ID
    pub fn agent_id(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }

    // NFT ID of collection
    pub fn collection(&self) -> ScImmutableNftID {
        ScImmutableNftID::new(self.proxy.root(PARAM_COLLECTION))
    }
}

#[derive(Clone)]
pub struct MutableAccountNFTAmountInCollectionParams {
    pub(crate) proxy: Proxy,
}

impl MutableAccountNFTAmountInCollectionParams {
    // account agent ID
    pub fn agent_id(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }

    // NFT ID of collection
    pub fn collection(&self) -> ScMutableNftID {
        ScMutableNftID::new(self.proxy.root(PARAM_COLLECTION))
    }
}

#[derive(Clone)]
pub struct ImmutableAccountNFTsParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableAccountNFTsParams {
    pub fn new() -> ImmutableAccountNFTsParams {
        ImmutableAccountNFTsParams {
            proxy: params_proxy(),
        }
    }

    // account agent ID
    pub fn agent_id(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct MutableAccountNFTsParams {
    pub(crate) proxy: Proxy,
}

impl MutableAccountNFTsParams {
    // account agent ID
    pub fn agent_id(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct ImmutableAccountNFTsInCollectionParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableAccountNFTsInCollectionParams {
    pub fn new() -> ImmutableAccountNFTsInCollectionParams {
        ImmutableAccountNFTsInCollectionParams {
            proxy: params_proxy(),
        }
    }

    // account agent ID
    pub fn agent_id(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }

    // NFT ID of collection
    pub fn collection(&self) -> ScImmutableNftID {
        ScImmutableNftID::new(self.proxy.root(PARAM_COLLECTION))
    }
}

#[derive(Clone)]
pub struct MutableAccountNFTsInCollectionParams {
    pub(crate) proxy: Proxy,
}

impl MutableAccountNFTsInCollectionParams {
    // account agent ID
    pub fn agent_id(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }

    // NFT ID of collection
    pub fn collection(&self) -> ScMutableNftID {
        ScMutableNftID::new(self.proxy.root(PARAM_COLLECTION))
    }
}

#[derive(Clone)]
pub struct ImmutableBalanceParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableBalanceParams {
    pub fn new() -> ImmutableBalanceParams {
        ImmutableBalanceParams {
            proxy: params_proxy(),
        }
    }

    // account agent ID
    pub fn agent_id(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct MutableBalanceParams {
    pub(crate) proxy: Proxy,
}

impl MutableBalanceParams {
    // account agent ID
    pub fn agent_id(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct ImmutableBalanceBaseTokenParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableBalanceBaseTokenParams {
    pub fn new() -> ImmutableBalanceBaseTokenParams {
        ImmutableBalanceBaseTokenParams {
            proxy: params_proxy(),
        }
    }

    // account agent ID
    pub fn agent_id(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct MutableBalanceBaseTokenParams {
    pub(crate) proxy: Proxy,
}

impl MutableBalanceBaseTokenParams {
    // account agent ID
    pub fn agent_id(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct ImmutableBalanceNativeTokenParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableBalanceNativeTokenParams {
    pub fn new() -> ImmutableBalanceNativeTokenParams {
        ImmutableBalanceNativeTokenParams {
            proxy: params_proxy(),
        }
    }

    // account agent ID
    pub fn agent_id(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }

    // native token ID
    pub fn token_id(&self) -> ScImmutableTokenID {
        ScImmutableTokenID::new(self.proxy.root(PARAM_TOKEN_ID))
    }
}

#[derive(Clone)]
pub struct MutableBalanceNativeTokenParams {
    pub(crate) proxy: Proxy,
}

impl MutableBalanceNativeTokenParams {
    // account agent ID
    pub fn agent_id(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }

    // native token ID
    pub fn token_id(&self) -> ScMutableTokenID {
        ScMutableTokenID::new(self.proxy.root(PARAM_TOKEN_ID))
    }
}

#[derive(Clone)]
pub struct ImmutableGetAccountNonceParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableGetAccountNonceParams {
    pub fn new() -> ImmutableGetAccountNonceParams {
        ImmutableGetAccountNonceParams {
            proxy: params_proxy(),
        }
    }

    // account agent ID
    pub fn agent_id(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct MutableGetAccountNonceParams {
    pub(crate) proxy: Proxy,
}

impl MutableGetAccountNonceParams {
    // account agent ID
    pub fn agent_id(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(PARAM_AGENT_ID))
    }
}

#[derive(Clone)]
pub struct ImmutableNativeTokenParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableNativeTokenParams {
    pub fn new() -> ImmutableNativeTokenParams {
        ImmutableNativeTokenParams {
            proxy: params_proxy(),
        }
    }

    // serial number of the foundry
    pub fn foundry_sn(&self) -> ScImmutableUint32 {
        ScImmutableUint32::new(self.proxy.root(PARAM_FOUNDRY_SN))
    }
}

#[derive(Clone)]
pub struct MutableNativeTokenParams {
    pub(crate) proxy: Proxy,
}

impl MutableNativeTokenParams {
    // serial number of the foundry
    pub fn foundry_sn(&self) -> ScMutableUint32 {
        ScMutableUint32::new(self.proxy.root(PARAM_FOUNDRY_SN))
    }
}

#[derive(Clone)]
pub struct ImmutableNftDataParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableNftDataParams {
    pub fn new() -> ImmutableNftDataParams {
        ImmutableNftDataParams {
            proxy: params_proxy(),
        }
    }

    // NFT ID
    pub fn nft_id(&self) -> ScImmutableNftID {
        ScImmutableNftID::new(self.proxy.root(PARAM_NFT_ID))
    }
}

#[derive(Clone)]
pub struct MutableNftDataParams {
    pub(crate) proxy: Proxy,
}

impl MutableNftDataParams {
    // NFT ID
    pub fn nft_id(&self) -> ScMutableNftID {
        ScMutableNftID::new(self.proxy.root(PARAM_NFT_ID))
    }
}
