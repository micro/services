import * as m3o from "@m3o/m3o-node";

export class ForexService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Returns the data for the previous close
  history(request: HistoryRequest): Promise<HistoryResponse> {
    return this.client.call(
      "forex",
      "History",
      request
    ) as Promise<HistoryResponse>;
  }
  // Get the latest price for a given forex ticker
  price(request: PriceRequest): Promise<PriceResponse> {
    return this.client.call(
      "forex",
      "Price",
      request
    ) as Promise<PriceResponse>;
  }
  // Get the latest quote for the forex
  quote(request: QuoteRequest): Promise<QuoteResponse> {
    return this.client.call(
      "forex",
      "Quote",
      request
    ) as Promise<QuoteResponse>;
  }
}

export interface HistoryRequest {
  // the forex symbol e.g GBPUSD
  symbol?: string;
}

export interface HistoryResponse {
  // the close price
  close?: number;
  // the date
  date?: string;
  // the peak price
  high?: number;
  // the low price
  low?: number;
  // the open price
  open?: number;
  // the forex symbol
  symbol?: string;
  // the volume
  volume?: number;
}

export interface PriceRequest {
  // forex symbol e.g GBPUSD
  symbol?: string;
}

export interface PriceResponse {
  // the last price
  price?: number;
  // the forex symbol e.g GBPUSD
  symbol?: string;
}

export interface QuoteRequest {
  // the forex symbol e.g GBPUSD
  symbol?: string;
}

export interface QuoteResponse {
  // the asking price
  askPrice?: number;
  // the bidding price
  bidPrice?: number;
  // the forex symbol
  symbol?: string;
  // the UTC timestamp of the quote
  timestamp?: string;
}
