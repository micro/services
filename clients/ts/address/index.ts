import * as m3o from "@m3o/m3o-node";

export class AddressService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Lookup a list of UK addresses by postcode
  lookupPostcode(
    request: LookupPostcodeRequest
  ): Promise<LookupPostcodeResponse> {
    return this.client.call(
      "address",
      "LookupPostcode",
      request
    ) as Promise<LookupPostcodeResponse>;
  }
}

export interface LookupPostcodeRequest {
  // UK postcode e.g SW1A 2AA
  postcode?: string;
}

export interface LookupPostcodeResponse {
  addresses?: Record[];
}

export interface Record {
  // building name
  buildingName?: string;
  // the county
  county?: string;
  // line one of address
  lineOne?: string;
  // line two of address
  lineTwo?: string;
  // dependent locality
  locality?: string;
  // organisation if present
  organisation?: string;
  // the postcode
  postcode?: string;
  // the premise
  premise?: string;
  // street name
  street?: string;
  // the complete address
  summary?: string;
  // post town
  town?: string;
}
