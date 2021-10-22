import * as m3o from "@m3o/m3o-node";

export class EmojiService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Find an emoji by its alias e.g :beer:
  find(request: FindRequest): Promise<FindResponse> {
    return this.client.call("emoji", "Find", request) as Promise<FindResponse>;
  }
  // Get the flag for a country. Requires country code e.g GB for great britain
  flag(request: FlagRequest): Promise<FlagResponse> {
    return this.client.call("emoji", "Flag", request) as Promise<FlagResponse>;
  }
  // Print text and renders the emojis with aliases e.g
  // let's grab a :beer: becomes let's grab a 🍺
  print(request: PrintRequest): Promise<PrintResponse> {
    return this.client.call(
      "emoji",
      "Print",
      request
    ) as Promise<PrintResponse>;
  }
  // Send an emoji to anyone via SMS. Messages are sent in the form '<message> Sent from <from>'
  send(request: SendRequest): Promise<SendResponse> {
    return this.client.call("emoji", "Send", request) as Promise<SendResponse>;
  }
}

export interface FindRequest {
  // the alias code e.g :beer:
  alias?: string;
}

export interface FindResponse {
  // the unicode emoji 🍺
  emoji?: string;
}

export interface FlagRequest {
  // country code e.g GB
  code?: string;
}

export interface FlagResponse {
  // the emoji flag
  flag?: string;
}

export interface PrintRequest {
  // text including any alias e.g let's grab a :beer:
  text?: string;
}

export interface PrintResponse {
  // text with rendered emojis
  text?: string;
}

export interface SendRequest {
  // the name of the sender from e.g Alice
  from?: string;
  // message to send including emoji aliases
  message?: string;
  // phone number to send to (including international dialing code)
  to?: string;
}

export interface SendResponse {
  // whether or not it succeeded
  success?: boolean;
}
