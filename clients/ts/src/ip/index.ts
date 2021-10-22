import * as m3o from "@m3o/m3o-node";

export class IpService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Lookup the geolocation information for an IP address
  lookup(request: LookupRequest): Promise<LookupResponse> {
    return this.client.call("ip", "Lookup", request) as Promise<LookupResponse>;
  }
}

export interface LookupRequest {
  // IP to lookup
  ip?: string;
}

export interface LookupResponse {
  // Autonomous system number
  asn?: number;
  // Name of the city
  city?: string;
  // Name of the continent
  continent?: string;
  // Name of the country
  country?: string;
  // IP of the query
  ip?: string;
  // Latitude e.g 52.523219
  latitude?: number;
  // Longitude e.g 13.428555
  longitude?: number;
  // Timezone e.g Europe/Rome
  timezone?: string;
}
