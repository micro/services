import * as m3o from "@m3o/m3o-node";

export class AnswerService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Ask a question and receive an instant answer
  question(request: QuestionRequest): Promise<QuestionResponse> {
    return this.client.call(
      "answer",
      "Question",
      request
    ) as Promise<QuestionResponse>;
  }
}

export interface QuestionRequest {
  // the question to answer
  query?: string;
}

export interface QuestionResponse {
  // the answer to your question
  answer?: string;
  // any related image
  image?: string;
  // a related url
  url?: string;
}
