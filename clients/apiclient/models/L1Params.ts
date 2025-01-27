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

import { BaseToken } from '../models/BaseToken';
import { ProtocolParameters } from '../models/ProtocolParameters';
import { HttpFile } from '../http/http';

export class L1Params {
    'baseToken': BaseToken;
    /**
    * The max payload size
    */
    'maxPayloadSize': number;
    'protocol': ProtocolParameters;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "baseToken",
            "baseName": "baseToken",
            "type": "BaseToken",
            "format": ""
        },
        {
            "name": "maxPayloadSize",
            "baseName": "maxPayloadSize",
            "type": "number",
            "format": "int32"
        },
        {
            "name": "protocol",
            "baseName": "protocol",
            "type": "ProtocolParameters",
            "format": ""
        }    ];

    static getAttributeTypeMap() {
        return L1Params.attributeTypeMap;
    }

    public constructor() {
    }
}

