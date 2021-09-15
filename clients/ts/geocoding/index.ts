import * as m3o from "@m3o/m3o-node";

export class GeocodingService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Lookup returns a geocoded address including normalized address and gps coordinates. All fields are optional, provide more to get more accurate results
  lookup(request: LookupRequest): Promise<LookupResponse> {
    return this.client.call(
      "geocoding",
      "Lookup",
      request
    ) as Promise<LookupResponse>;
  }
  // Reverse lookup an address from gps coordinates
  reverse(request: ReverseRequest): Promise<ReverseResponse> {
    return this.client.call(
      "geocoding",
      "Reverse",
      request
    ) as Promise<ReverseResponse>;
  }
}

export interface Address {
  city?: string;
  country?: string;
  lineOne?: string;
  lineTwo?: string;
  postcode?: string;
}

export interface Location {
  latitude?: number;
  longitude?: number;
}

export interface LookupRequest {
  address?: string;
  city?: string;
  country?: string;
  postcode?: string;
}

export interface LookupResponse {
  address?: { [key: string]: any };
  location?: { [key: string]: any };
}

export interface ReverseRequest {
  latitude?: number;
  longitude?: number;
}

export interface ReverseResponse {
  address?: { [key: string]: any };
  location?: { [key: string]: any };
}
