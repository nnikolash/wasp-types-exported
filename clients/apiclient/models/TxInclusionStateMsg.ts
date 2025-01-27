/**
 * Wasp API
 * REST API for the Wasp node
 *
 * OpenAPI spec version: 0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { HttpFile } from '../http/http';

export class TxInclusionStateMsg {
    /**
    * The inclusion state
    */
    'state': string;
    /**
    * The transaction ID
    */
    'txId': string;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "state",
            "baseName": "state",
            "type": "string",
            "format": "string"
        },
        {
            "name": "txId",
            "baseName": "txId",
            "type": "string",
            "format": "string"
        }    ];

    static getAttributeTypeMap() {
        return TxInclusionStateMsg.attributeTypeMap;
    }

    public constructor() {
    }
}

