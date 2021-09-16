import * as m3o from "@m3o/m3o-node";

export class SentimentService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Analyze and score a piece of text
  analyze(request: AnalyzeRequest): Promise<AnalyzeResponse> {
    return this.client.call(
      "sentiment",
      "Analyze",
      request
    ) as Promise<AnalyzeResponse>;
  }
}

export interface AnalyzeRequest {
  // The language. Defaults to english.
  lang?: string;
  // The text to analyze
  text?: string;
}

export interface AnalyzeResponse {
  // The score of the text {positive is 1, negative is 0}
  score?: number;
}
