import * as m3o from "@m3o/m3o-node";

export class HolidaysService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  //
  countries(request: CountriesRequest): Promise<CountriesResponse> {
    return this.client.call(
      "holidays",
      "Countries",
      request
    ) as Promise<CountriesResponse>;
  }
  //
  list(request: ListRequest): Promise<ListResponse> {
    return this.client.call(
      "holidays",
      "List",
      request
    ) as Promise<ListResponse>;
  }
}

export interface CountriesRequest {}

export interface CountriesResponse {
  countries?: Country[];
}

export interface Country {
  // The 2 letter country code (as defined in ISO 3166-1 alpha-2)
  code?: string;
  // The English name of the country
  name?: string;
}

export interface Holiday {
  // the country this holiday occurs in
  countryCode?: string;
  // date of the holiday in yyyy-mm-dd format
  date?: string;
  // the local name of the holiday
  localName?: string;
  // the name of the holiday in English
  name?: string;
  // the regions within the country that observe this holiday (if not all of them)
  regions?: string[];
  // the type of holiday Public, Bank, School, Authorities, Optional, Observance
  types?: string[];
}

export interface ListRequest {
  // The 2 letter country code (as defined in ISO 3166-1 alpha-2)
  countryCode?: string;
  // The year to list holidays for
  year?: number;
}

export interface ListResponse {
  holidays?: Holiday[];
}
