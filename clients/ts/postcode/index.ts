import * as m3o from "@m3o/m3o-node";

export class PostcodeService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Lookup a postcode to retrieve the related region, county, etc
  lookup(request: LookupRequest): Promise<LookupResponse> {
    return this.client.call(
      "postcode",
      "Lookup",
      request
    ) as Promise<LookupResponse>;
  }
  // Return a random postcode and its related info
  random(request: RandomRequest): Promise<RandomResponse> {
    return this.client.call(
      "postcode",
      "Random",
      request
    ) as Promise<RandomResponse>;
  }
  // Validate a postcode.
  validate(request: ValidateRequest): Promise<ValidateResponse> {
    return this.client.call(
      "postcode",
      "Validate",
      request
    ) as Promise<ValidateResponse>;
  }
}

export interface LookupRequest {
  // UK postcode e.g SW1A 2AA
  postcode?: string;
}

export interface LookupResponse {
  // country e.g United Kingdom
  country?: string;
  // e.g Westminster
  district?: string;
  // e.g 51.50354
  latitude?: number;
  // e.g -0.127695
  longitude?: number;
  // UK postcode e.g SW1A 2AA
  postcode?: string;
  // related region e.g London
  region?: string;
  // e.g St James's
  ward?: string;
}

export interface RandomRequest {}

export interface RandomResponse {
  // country e.g United Kingdom
  country?: string;
  // e.g Westminster
  district?: string;
  // e.g 51.50354
  latitude?: number;
  // e.g -0.127695
  longitude?: number;
  // UK postcode e.g SW1A 2AA
  postcode?: string;
  // related region e.g London
  region?: string;
  // e.g St James's
  ward?: string;
}

export interface ValidateRequest {
  // UK postcode e.g SW1A 2AA
  postcode?: string;
}

export interface ValidateResponse {
  // Is the postcode valid (true) or not (false)
  valid?: boolean;
}
