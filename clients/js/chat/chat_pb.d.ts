import * as jspb from "google-protobuf"

export class NewRequest extends jspb.Message {
  getUserIdsList(): Array<string>;
  setUserIdsList(value: Array<string>): void;
  clearUserIdsList(): void;
  addUserIds(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewRequest.AsObject;
  static toObject(includeInstance: boolean, msg: NewRequest): NewRequest.AsObject;
  static serializeBinaryToWriter(message: NewRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewRequest;
  static deserializeBinaryFromReader(message: NewRequest, reader: jspb.BinaryReader): NewRequest;
}

export namespace NewRequest {
  export type AsObject = {
    userIdsList: Array<string>,
  }
}

export class NewResponse extends jspb.Message {
  getChatId(): string;
  setChatId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewResponse.AsObject;
  static toObject(includeInstance: boolean, msg: NewResponse): NewResponse.AsObject;
  static serializeBinaryToWriter(message: NewResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewResponse;
  static deserializeBinaryFromReader(message: NewResponse, reader: jspb.BinaryReader): NewResponse;
}

export namespace NewResponse {
  export type AsObject = {
    chatId: string,
  }
}

export class HistoryRequest extends jspb.Message {
  getChatId(): string;
  setChatId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HistoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: HistoryRequest): HistoryRequest.AsObject;
  static serializeBinaryToWriter(message: HistoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HistoryRequest;
  static deserializeBinaryFromReader(message: HistoryRequest, reader: jspb.BinaryReader): HistoryRequest;
}

export namespace HistoryRequest {
  export type AsObject = {
    chatId: string,
  }
}

export class HistoryResponse extends jspb.Message {
  getMessagesList(): Array<Message>;
  setMessagesList(value: Array<Message>): void;
  clearMessagesList(): void;
  addMessages(value?: Message, index?: number): Message;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HistoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: HistoryResponse): HistoryResponse.AsObject;
  static serializeBinaryToWriter(message: HistoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HistoryResponse;
  static deserializeBinaryFromReader(message: HistoryResponse, reader: jspb.BinaryReader): HistoryResponse;
}

export namespace HistoryResponse {
  export type AsObject = {
    messagesList: Array<Message.AsObject>,
  }
}

export class SendRequest extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): void;

  getChatId(): string;
  setChatId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getSubject(): string;
  setSubject(value: string): void;

  getText(): string;
  setText(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SendRequest): SendRequest.AsObject;
  static serializeBinaryToWriter(message: SendRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendRequest;
  static deserializeBinaryFromReader(message: SendRequest, reader: jspb.BinaryReader): SendRequest;
}

export namespace SendRequest {
  export type AsObject = {
    clientId: string,
    chatId: string,
    userId: string,
    subject: string,
    text: string,
  }
}

export class SendResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SendResponse): SendResponse.AsObject;
  static serializeBinaryToWriter(message: SendResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendResponse;
  static deserializeBinaryFromReader(message: SendResponse, reader: jspb.BinaryReader): SendResponse;
}

export namespace SendResponse {
  export type AsObject = {
  }
}

export class Message extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getClientId(): string;
  setClientId(value: string): void;

  getChatId(): string;
  setChatId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getSentAt(): number;
  setSentAt(value: number): void;

  getSubject(): string;
  setSubject(value: string): void;

  getText(): string;
  setText(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Message.AsObject;
  static toObject(includeInstance: boolean, msg: Message): Message.AsObject;
  static serializeBinaryToWriter(message: Message, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Message;
  static deserializeBinaryFromReader(message: Message, reader: jspb.BinaryReader): Message;
}

export namespace Message {
  export type AsObject = {
    id: string,
    clientId: string,
    chatId: string,
    userId: string,
    sentAt: number,
    subject: string,
    text: string,
  }
}

