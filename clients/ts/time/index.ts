import * as m3o from "@m3o/m3o-node";

export class TimeService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Get the current time
  now(request: NowRequest): Promise<NowResponse> {
    return this.client.call("time", "Now", request) as Promise<NowResponse>;
  }
  // Get the timezone info for a specific location
  zone(request: ZoneRequest): Promise<ZoneResponse> {
    return this.client.call("time", "Zone", request) as Promise<ZoneResponse>;
  }
}

export interface NowRequest {
  // optional location, otherwise returns UTC
  location?: string;
}

export interface NowResponse {
  // the current time as HH:MM:SS
  localtime?: string;
  // the location as Europe/London
  location?: string;
  // timestamp as 2006-01-02T15:04:05.999999999Z07:00
  timestamp?: string;
  // the timezone as BST
  timezone?: string;
  // the unix timestamp
  unix?: number;
}

export interface ZoneRequest {
  // location to lookup e.g postcode, city, ip address
  location?: string;
}

export interface ZoneResponse {
  // the abbreviated code e.g BST
  abbreviation?: string;
  // country of the timezone
  country?: string;
  // is daylight savings
  dst?: boolean;
  // e.g 51.42
  latitude?: number;
  // the local time
  localtime?: string;
  // location requested
  location?: string;
  // e.g -0.37
  longitude?: number;
  // region of timezone
  region?: string;
  // the timezone e.g Europe/London
  timezone?: string;
}
