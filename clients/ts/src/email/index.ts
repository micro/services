import * as m3o from "@m3o/m3o-node";

export class EmailService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Send an email by passing in from, to, subject, and a text or html body
  send(request: SendRequest): Promise<SendResponse> {
    return this.client.call("email", "Send", request) as Promise<SendResponse>;
  }
}

export interface SendRequest {
  // the display name of the sender
  from?: string;
  // the html body
  htmlBody?: string;
  // an optional reply to email address
  replyTo?: string;
  // the email subject
  subject?: string;
  // the text body
  textBody?: string;
  // the email address of the recipient
  to?: string;
}

export interface SendResponse {}
