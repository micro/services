import * as m3o from "@m3o/m3o-node";

export class HelloworldService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Call returns a personalised "Hello $name" response
  call(request: CallRequest): Promise<CallResponse> {
    return this.client.call(
      "helloworld",
      "Call",
      request
    ) as Promise<CallResponse>;
  }
  // Stream returns a stream of "Hello $name" responses
  stream(request: StreamRequest): Promise<StreamResponse> {
    return this.client.call(
      "helloworld",
      "Stream",
      request
    ) as Promise<StreamResponse>;
  }
}

export interface CallRequest {
  name?: string;
}

export interface CallResponse {
  message?: string;
}

export interface StreamRequest {
  // the number of messages to send back
  messages?: number;
  name?: string;
}

export interface StreamResponse {
  message?: string;
}
