import * as m3o from "@m3o/m3o-node";

export class CurrencyService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Codes returns the supported currency codes for the API
  codes(request: CodesRequest): Promise<CodesResponse> {
    return this.client.call(
      "currency",
      "Codes",
      request
    ) as Promise<CodesResponse>;
  }
  // Convert returns the currency conversion rate between two pairs e.g USD/GBP
  convert(request: ConvertRequest): Promise<ConvertResponse> {
    return this.client.call(
      "currency",
      "Convert",
      request
    ) as Promise<ConvertResponse>;
  }
  // Returns the historic rates for a currency on a given date
  history(request: HistoryRequest): Promise<HistoryResponse> {
    return this.client.call(
      "currency",
      "History",
      request
    ) as Promise<HistoryResponse>;
  }
  // Rates returns the currency rates for a given code e.g USD
  rates(request: RatesRequest): Promise<RatesResponse> {
    return this.client.call(
      "currency",
      "Rates",
      request
    ) as Promise<RatesResponse>;
  }
}

export interface Code {
  // e.g United States Dollar
  currency?: string;
  // e.g USD
  name?: string;
}

export interface CodesRequest {}

export interface CodesResponse {
  codes?: Code[];
}

export interface ConvertRequest {
  // optional amount to convert e.g 10.0
  amount?: number;
  // base code to convert from e.g USD
  from?: string;
  // target code to convert to e.g GBP
  to?: string;
}

export interface ConvertResponse {
  // converted amount e.g 7.10
  amount?: number;
  // the base code e.g USD
  from?: string;
  // conversion rate e.g 0.71
  rate?: number;
  // the target code e.g GBP
  to?: string;
}

export interface HistoryRequest {
  // currency code e.g USD
  code?: string;
  // date formatted as YYYY-MM-DD
  date?: string;
}

export interface HistoryResponse {
  // The code of the request
  code?: string;
  // The date requested
  date?: string;
  // The rate for the day as code:rate
  rates?: { [key: string]: number };
}

export interface RatesRequest {
  // The currency code to get rates for e.g USD
  code?: string;
}

export interface RatesResponse {
  // The code requested e.g USD
  code?: string;
  // The rates for the given code as key-value pairs code:rate
  rates?: { [key: string]: number };
}
