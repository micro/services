import * as m3o from "@m3o/m3o-node";

export class StockService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Get the historic open-close for a given day
  history(request: HistoryRequest): Promise<HistoryResponse> {
    return this.client.call(
      "stock",
      "History",
      request
    ) as Promise<HistoryResponse>;
  }
  // Get the historic order book and each trade by timestamp
  orderBook(request: OrderBookRequest): Promise<OrderBookResponse> {
    return this.client.call(
      "stock",
      "OrderBook",
      request
    ) as Promise<OrderBookResponse>;
  }
  // Get the last price for a given stock ticker
  price(request: PriceRequest): Promise<PriceResponse> {
    return this.client.call(
      "stock",
      "Price",
      request
    ) as Promise<PriceResponse>;
  }
  // Get the last quote for the stock
  quote(request: QuoteRequest): Promise<QuoteResponse> {
    return this.client.call(
      "stock",
      "Quote",
      request
    ) as Promise<QuoteResponse>;
  }
}

export interface HistoryRequest {
  // date to retrieve as YYYY-MM-DD
  date?: string;
  // the stock symbol e.g AAPL
  stock?: string;
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
  // the stock symbol
  symbol?: string;
  // the volume
  volume?: number;
}

export interface Order {
  // the asking price
  askPrice?: number;
  // the ask size
  askSize?: number;
  // the bidding price
  bidPrice?: number;
  // the bid size
  bidSize?: number;
  // the UTC timestamp of the quote
  timestamp?: string;
}

export interface OrderBookRequest {
  // the date in format YYYY-MM-dd
  date?: string;
  // optional RFC3339Nano end time e.g 2006-01-02T15:04:05.999999999Z07:00
  end?: string;
  // limit number of prices
  limit?: number;
  // optional RFC3339Nano start time e.g 2006-01-02T15:04:05.999999999Z07:00
  start?: string;
  // stock to retrieve e.g AAPL
  stock?: string;
}

export interface OrderBookResponse {
  // date of the request
  date?: string;
  // list of orders
  orders?: Order[];
  // the stock symbol
  symbol?: string;
}

export interface PriceRequest {
  // stock symbol e.g AAPL
  symbol?: string;
}

export interface PriceResponse {
  // the last price
  price?: number;
  // the stock symbol e.g AAPL
  symbol?: string;
}

export interface QuoteRequest {
  // the stock symbol e.g AAPL
  symbol?: string;
}

export interface QuoteResponse {
  // the asking price
  askPrice?: number;
  // the ask size
  askSize?: number;
  // the bidding price
  bidPrice?: number;
  // the bid size
  bidSize?: number;
  // the stock symbol
  symbol?: string;
  // the UTC timestamp of the quote
  timestamp?: string;
}
