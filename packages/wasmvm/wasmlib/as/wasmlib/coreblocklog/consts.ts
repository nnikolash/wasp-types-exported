// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import * as wasmtypes from '../wasmtypes';

export const ScName        = 'blocklog';
export const ScDescription = 'Block log contract';
export const HScName       = new wasmtypes.ScHname(0xf538ef2b);

export const ParamBlockIndex = 'n';
export const ParamRequestID  = 'u';

export const ResultBlockIndex       = 'n';
export const ResultBlockInfo        = 'i';
export const ResultEvent            = 'e';
export const ResultRequestID        = 'u';
export const ResultRequestIndex     = 'r';
export const ResultRequestProcessed = 'p';
export const ResultRequestReceipt   = 'd';
export const ResultRequestReceipts  = 'd';

export const ViewGetBlockInfo               = 'getBlockInfo';
export const ViewGetEventsForBlock          = 'getEventsForBlock';
export const ViewGetEventsForRequest        = 'getEventsForRequest';
export const ViewGetRequestIDsForBlock      = 'getRequestIDsForBlock';
export const ViewGetRequestReceipt          = 'getRequestReceipt';
export const ViewGetRequestReceiptsForBlock = 'getRequestReceiptsForBlock';
export const ViewIsRequestProcessed         = 'isRequestProcessed';

export const HViewGetBlockInfo               = new wasmtypes.ScHname(0xbe89f9b3);
export const HViewGetEventsForBlock          = new wasmtypes.ScHname(0x36232798);
export const HViewGetEventsForRequest        = new wasmtypes.ScHname(0x4f8d68e4);
export const HViewGetRequestIDsForBlock      = new wasmtypes.ScHname(0x5a20327a);
export const HViewGetRequestReceipt          = new wasmtypes.ScHname(0xb7f9534f);
export const HViewGetRequestReceiptsForBlock = new wasmtypes.ScHname(0x77e3beef);
export const HViewIsRequestProcessed         = new wasmtypes.ScHname(0xd57d50a9);
